package main

import (
	"context"
	"fmt"
	"github.com/akimdev15/melongo/playlist-server/spotify"
	"log"
	"net"

	"github.com/akimdev15/melongo/playlist-server/internal/database"
	"github.com/akimdev15/melongo/playlist-server/proto"
	"google.golang.org/grpc"
)

const gRPCPORT = "50002"

type PlaylistServer struct {
	proto.UnimplementedPlaylistServiceServer
	DB *database.Queries
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

	// 1. Get music from melon and search them
	//songs := mscraper.GetNewestSongsMelon("0900")
	//
	//var URIs []string
	//var searchResult spotify.TracksResponse
	//for _, song := range songs[:1] {
	//	artistInfo, err := spotify.SearchArtistID(song.Artist, req.AccessToken)
	//	// search artist id test
	//	if err != nil {
	//		fmt.Println("Error searching for artist ID. err: ", err)
	//		return nil, err
	//	}
	//	fmt.Printf("ArtistID: %v\n", artistInfo)
	//
	//	searchResult, err = spotify.SearchTrack(song.Title, artistInfo.Name, req.AccessToken)
	//	if err != nil || len(searchResult.Items.Tracks) <= 0 {
	//		fmt.Println("error while getting the search result. Error: ", err)
	//		continue
	//	}
	//
	//	URIs = append(URIs, searchResult.Items.Tracks[0].URI)
	//
	//	fmt.Printf("search result: %v\n", searchResult)
	//}
	//
	//addTrackRes, err := spotify.AddTrackToPlaylist(newPlaylistResponse.SpotifyPlaylistID, URIs, req.AccessToken)
	//
	//if err != nil {
	//	fmt.Println("Error adding tracks to the playlist")
	//}

}
