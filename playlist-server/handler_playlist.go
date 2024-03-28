package main

import (
	"fmt"
	"github.com/melongo/playlist-server/spotify"
	"net/http"
)

// AccessToken TODO - only for testing purpose. Should be REMOVED!!
const AccessToken = "BQA-_dXRuy4342MQBdqFd2VuOVDRW1EhC8MlJfu1A3MuBtIkdKUOCf4RGfn5_rOcJmgXUnoGMJZCIhDzAGrUUpiWIsdhBSRUwO-Nu5my1ZefyMjyQkwp0G-9k8x3Jfw0kbOlEtNZrwtfcDCQjTUGBvbQSG8L8ZTH2QgDRwAOA-4f-JTLK7oqIGvBZsipbBNEZr9Bypq0aKNhCT7HNaoK2gAlT7X0yEqlcauj7dnq0e2Wxo1c64uShqrQcp8jr_GItCAhpyAg97jv9Gz-Pm7F5AvSwQ"

func (apiCfg *apiConfig) testHandler(w http.ResponseWriter, r *http.Request) {
	// search artist id test
	artistID, err := spotify.SearchArtistID("아이유", AccessToken)
	if err != nil {
		fmt.Println("Error searching for artist ID. err: ", err)
		respondWithError(w, 401, fmt.Sprintf("Error getting the artist ID. err: %v\n", err))
	}
	fmt.Printf("ArtistID: %v\n", artistID)

	tracks, err := spotify.SearchTracksByArtist(artistID, AccessToken)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Error getting the tracks. err: %v\n", err))
	}
	for _, track := range tracks {
		fmt.Printf("URI: %s ->  Track Name: %s -> Artist Name: %s\n", track.URI, track.Name, track.Album.Artists[0].Name)
	}

	playlists, err := spotify.GetUserPlaylists(AccessToken)
	if err != nil {
		// return json error
		respondWithError(w, 401, fmt.Sprintf("Error getting the playlists. err: %v\n", err))
	}

	fmt.Println("Playlists: ", playlists)

	respondWithJSON(w, 200, tracks)

	// for testing
	//searchResult, err := apiCfg.searchTrack("홀씨", "아이유", AccessToken)
	//if err != nil {
	//	fmt.Println("error while getting the search result. Error: ", err)
	//}
	//
	//fmt.Printf("search result: %v\n", searchResult)
	//
	//item := searchResult.Tracks.Items[0]
	//
	//songName := item.Name
	//
	//// Artist name
	//artistName := item.Album.Artists[0].Name
	//
	//// Track URI
	//trackURI := item.URI
	//
	//fmt.Printf("Song Name: %s\n", songName)
	//fmt.Printf("Artist Name: %s\n", artistName)
	//fmt.Printf("Track URI: %s\n", trackURI)
}
