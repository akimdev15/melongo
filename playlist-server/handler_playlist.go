package main

import (
	"fmt"
	"github.com/akimdev15/melongo/playlist-server/spotify"
	"github.com/akimdev15/mscraper"
	"net/http"
)

// AccessToken TODO - only for testing purpose. Should be REMOVED!!
const AccessToken = "BQDTLi92KAD0GfTFEx2NgKx9_AJ98CRVOO4aCEuBT_4IiDm8dLZF71BHsL2-0-hqGGcsZF89fBXBmrmNJo_Livv-qksl7zv2qNqIs0xZQLeMYUumJn-gpY9R_4Oi3WXljEnG1V4RJLMghB8t4lvo7wKsqmiEHGyLoOrJuCxHqDR5NxmGSTUiCHtGXxW1-B5r9UPrL442Ji_1r7X2xonszCetTvgLW2_SIqSAnv2Esu1oQnl9qP7fW9rvKg224n-wX3_EmALsura-nOLiIM3TlmOECg"

func (apiCfg *apiConfig) testHandler(w http.ResponseWriter, r *http.Request) {

	songs := mscraper.GetNewestSongsMelon("0300")

	var searchResult spotify.TracksResponse
	for _, song := range songs[:1] {
		artistInfo, err := spotify.SearchArtistID(song.Artist, AccessToken)
		// search artist id test
		if err != nil {
			fmt.Println("Error searching for artist ID. err: ", err)
			respondWithError(w, 401, fmt.Sprintf("Error getting the artist ID. err: %v\n", err))
		}
		fmt.Printf("ArtistID: %v\n", artistInfo)

		searchResult, err = spotify.SearchTrack(song.Title, artistInfo.Name, AccessToken)
		if err != nil {
			fmt.Println("error while getting the search result. Error: ", err)
		}

		fmt.Printf("search result: %v\n", searchResult)
	}

	playlists, err := spotify.GetUserPlaylists(AccessToken)
	if err != nil {
		// return json error
		respondWithError(w, 401, fmt.Sprintf("Error getting the playlists. err: %v\n", err))
	}

	fmt.Println("Playlists: ", playlists)

	respondWithJSON(w, 200, searchResult)

}
