package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/akimdev15/melongo/playlist-server/spotify"
	"github.com/akimdev15/mscraper"

	"github.com/akimdev15/melongo/playlist-server/internal/database"
	"github.com/akimdev15/melongo/playlist-server/proto"
	"google.golang.org/grpc"
)

const gRPCPORT = "50002"

type PlaylistServer struct {
	proto.UnimplementedPlaylistServiceServer
	DB *database.Queries
}

// SongDB is a struct to store song information in the database
type SongDB struct {
	Rank   int32
	Title  string
	Artist string
	URI    string
	Date   time.Time
}

func (apiCfg *apiConfig) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPORT))
	if err != nil {
		log.Fatalf("Failed to listen for grpc %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterPlaylistServiceServer(grpcServer, &PlaylistServer{DB: apiCfg.DB})
	fmt.Printf("gRPC server start on PORT %s]\n", gRPCPORT)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC %v", err)
	}
}

func (playlistServer *PlaylistServer) CreatePlaylist(ctx context.Context, req *proto.CreatePlaylistRequest) (*proto.CreatePlaylistResponse, error) {
	fmt.Println("Received request to create playlist using gRPC")
	// 0. Create a new playlist
	newPlaylistResponse, err := spotify.CreateNewPlaylist(req.PlaylistName, req.Description, req.IsPublic, req.UserID, req.AccessToken)
	if err != nil {
		fmt.Println("Error creating new playlist")
		return nil, err
	}

	fmt.Println("new playlist: ", newPlaylistResponse)

	res := &proto.CreatePlaylistResponse{
		SpotifyPlaylistID: newPlaylistResponse.SpotifyPlaylistID,
		ExternalURL:       newPlaylistResponse.URI,
		Name:              newPlaylistResponse.Name,
	}
	return res, nil

}

func (PlaylistServer *PlaylistServer) CreateMelonTop100(ctx context.Context, req *proto.CreateMelonTop100Request) (*proto.CreateMelonTop100Response, error) {
	// TODO - need to use cache (takes too long)
	songs := mscraper.GetMelonTop100Songs()

	var wg sync.WaitGroup
	uriChan := make(chan string, len(songs))

	for _, song := range songs {
		wg.Add(1)
		go func(song mscraper.Song) {
			defer wg.Done()
			track, err := spotify.SearchTrack(song.Title, song.Artist, req.AccessToken)
			if err != nil {
				// TODO - should collect the missed tracks in chan and add these info to DB or something
				fmt.Printf("Nothing found for song: %s and artist: %s\n", song.Title, song.Artist)
				return
			}
			if track != nil && track.URI != "" {
				uriChan <- track.URI
			}
		}(song)
	}

	wg.Wait()
	close(uriChan)

	// Collect URIs from the channel
	var uris []string
	for uri := range uriChan {
		uris = append(uris, uri)
	}

	// Return the response before adding tracks to the playlist
	response := &proto.CreateMelonTop100Response{
		Status: fmt.Sprintf("Added %d tracks and missed %d tracks", len(uris), len(songs)-len(uris)),
	}

	// Add tracks to the playlist after sending the response to reduce time
	go func() {
		_, err := spotify.AddTrackToPlaylist(req.PlaylistID, uris, req.AccessToken)
		if err != nil {
			fmt.Println("Error adding tracks to playlist: ", err)
		}
	}()

	return response, nil
}

func (PlaylistServer *PlaylistServer) SaveMelonTop100DB(ctx context.Context, req *proto.SaveMelonTop100DBRequest) (*proto.SaveMelonTop100DBResponse, error) {
	// get today's date in the format of "YYYY-MM-DD"
	date := getKST()

	go PlaylistServer.searchTracksAndSaveToDB(date, req.AccessToken)

	response := &proto.SaveMelonTop100DBResponse{
		Status: fmt.Sprintf("Saving top 100 melon tracks for the date: %s", date),
	}

	return response, nil
}

// ------------------ Helper Functions ------------------

func (PlaylistServer *PlaylistServer) searchTracksAndSaveToDB(date time.Time, accessToken string) {
	songs := mscraper.GetMelonTop100Songs()

	var wg sync.WaitGroup
	songChan := make(chan SongDB, len(songs))

	// Make the db insert non-blocking
	// TODO - Find a way to do bulk insert
	go PlaylistServer.saveTrackToDB(songChan)

	// Search tracks (this gets executed before the db insert above)
	for i, song := range songs {
		wg.Add(1)
		go PlaylistServer.processSong(song, i, date, accessToken, songChan, &wg)
	}

	wg.Wait()
	close(songChan)
}

func (playlistServer *PlaylistServer) saveTrackToDB(songChan <-chan SongDB) {
	for songDB := range songChan {
		_, err := playlistServer.DB.CreateTrack(context.Background(), database.CreateTrackParams{
			Rank:   songDB.Rank,
			Title:  songDB.Title,
			Artist: songDB.Artist,
			Uri:    songDB.URI,
			Date:   songDB.Date,
		})

		if err != nil {
			fmt.Printf("Error saving song to DB: %v. Error: %v", songDB, err)
			return
		}
	}
}

func (playlistServer *PlaylistServer) processSong(song mscraper.Song, index int, date time.Time, accessToken string, songChan chan<- SongDB, wg *sync.WaitGroup) {
	defer wg.Done()

	track, err := spotify.SearchTrack(song.Title, song.Artist, accessToken)
	if err != nil {
		playlistServer.handleTrackSearchError(song, index, date, songChan)
		return
	}

	// If track successfully found from Spotify, add it to the songChan
	if track != nil && track.URI != "" {
		songChan <- SongDB{
			Rank:   int32(index + 1),
			Title:  track.Name,
			Artist: track.Artist,
			URI:    track.URI,
			Date:   date,
		}
	}
}

func (playlistServer *PlaylistServer) handleTrackSearchError(song mscraper.Song, index int, date time.Time, songChan chan<- SongDB) {
	resolvedTrack, err := playlistServer.DB.GetResolvedTrack(context.Background(), database.GetResolvedTrackParams{
		MissedTitle:  song.Title,
		MissedArtist: song.Artist,
	})

	if err != nil {
		fmt.Printf("Error getting resolved track: %v. Adding it to the missed track DB.\n", err)
		_, err := playlistServer.DB.CreateMissedTrack(context.Background(), database.CreateMissedTrackParams{
			Rank:   int32(index + 1),
			Title:  song.Title,
			Artist: song.Artist,
			Date:   date,
		})
		if err != nil {
			fmt.Printf("Error saving missed track to DB: %v. Error: %v\n", song, err)
		}
		return
	}

	// Resolved track found. Add it to the songChan
	songChan <- SongDB{
		Rank:   int32(index + 1),
		Title:  resolvedTrack.Title,
		Artist: resolvedTrack.Artist,
		URI:    resolvedTrack.Uri,
		Date:   date,
	}
}

// getKST returns the current date in KST timezone
// do date.Format("2006-01-02") to get the date in the format of "YYYY-MM-DD"
func getKST() time.Time {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}
	}
	return time.Now().In(loc)
}
