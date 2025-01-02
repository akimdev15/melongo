package main

import (
	"fmt"
	"net/http"

	"github.com/akimdev15/melongo/playlist-server/spotify"
	"github.com/akimdev15/mscraper"
)

// AccessToken TODO - only for testing purpose. Should be REMOVED!!
const AccessToken = ""

func getSongs() []mscraper.Song {
	songs := mscraper.GetNewestSongsMelon("0300")
	return songs
}

func (apiCfg *apiConfig) testHandler(w http.ResponseWriter, r *http.Request) {
	songs := getSongs()
	uris := []string{}
	for i, song := range songs {
		if i > 5 {
			break
		}
		fmt.Printf("Searching for song: %s artist: %s\n", song.Title, song.Artist)
		track, err := spotify.SearchTrack(song.Title, song.Artist, AccessToken)
		if err != nil {
			fmt.Printf("Nothing found for song: %s and artist: %s\n", song.Title, song.Artist)
		}
		fmt.Println("Found track: ", track)
		if track != nil && track.URI != "" {
			uris = append(uris, track.URI)
		}
	}

	spotify.AddTrackToPlaylist("0msfdSZz5ZKXibCW6uZlvU", uris, AccessToken)

	respondWithJSON(w, 200, uris)
}

func (apiCfg *apiConfig) testNewAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	// albums := mscraper.GetNewestAlbumFromMelon()
	// tracks, err := spotify.SearchTracksFromAlbum(albums[0].Name, albums[0].Artist, AccessToken)

	// This one works
	tracks, err := spotify.SearchTracksFromAlbum("아무렇지 않게", "DK", AccessToken)

	// This works but artis name is empty
	// tracks, err := spotify.SearchTracksFromAlbum("#2024: 가장자리", "Minit", AccessToken)

	if err != nil {
		fmt.Println("Error searching tracks: ", err)
		return
	}

	uris := []string{}
	for _, track := range tracks {
		uris = append(uris, track.URI)
	}

	// spotify.AddTrackToPlaylist("0msfdSZz5ZKXibCW6uZlvU", uris, AccessToken)

	respondWithJSON(w, 200, tracks)
}
