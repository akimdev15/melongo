package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akimdev15/melongo/internal/database"
)

// handleCreatePlaylist creates a new playlist
func (apiCfg *apiConfig) handleCreatePlaylist(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON%s: ", err))
		return
	}

	playlist, err := apiCfg.DB.CreatePlaylist(r.Context(), database.CreatePlaylistParams{
		Name:   params.Name,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating a playlist: %s. Error:  %s\n", params.Name, err))
		return
	}

	respondWithJSON(w, 201, databasePlaylistToPlaylist(playlist))
}

// handleGetAllPlaylists returns all the playlists for the user making the request
func (apiCfg *apiConfig) handleGetAllPlaylists(w http.ResponseWriter, r *http.Request, user database.User) {
	playlists, err := apiCfg.DB.GetPlaylistsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting playlists for the user: %s. Error:  %s\n", user.Name, err))
		return
	}

	respondWithJSON(w, 201, databasePlaylistsToPlaylists(playlists))
}

// handleGetPlaylist retrieves playlist by the playlist name and the username
func (apiCfg *apiConfig) handleGetPlaylist(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON%s: ", err))
		return
	}
	playlist, err := apiCfg.DB.GetPlaylist(r.Context(), database.GetPlaylistParams{
		UserID: user.ID,
		Name:   params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting playlist: %s for the user: %s. Error:  %s\n", params.Name, user.Name, err))
		return
	}
	respondWithJSON(w, 201, databasePlaylistToPlaylist(playlist))
}

// TODO - create handler to add music to the playlist
