package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akimdev15/melongo/internal/database"
	"github.com/akimdev15/mscraper"
	"github.com/google/uuid"
)

// handleDailyMusicScrape - save music to the database
func (apiCfg *apiConfig) handleDailyMusicScrape() {
	genres, err := apiCfg.DB.GetAllGenre(context.Background())
	if err != nil {
		fmt.Printf("Error getting the genre code from the DB. Err: %s", err)
	}
	for _, genre := range genres {
		songs := mscraper.GetNewestSongsMelon(genre.Code)
		for _, song := range songs {
			_, err := apiCfg.DB.CreateMusic(context.Background(), database.CreateMusicParams{
				ID:        uuid.New(),
				Name:      song.Title,
				Artist:    song.Artist,
				CreatedAt: time.Now().UTC(),
				Genre:     genre.ID,
			})
			if err != nil {
				fmt.Printf("Error creating the music: %v.\n  Error: %s", song, err)
			}
		}
		fmt.Printf("Saved %d songs for genre: %s", len(songs), genre.Name)
	}
}

// handleGetMusic - using query parameters. Either by title, artist or both
func (apiCfg *apiConfig) handleGetMusic(w http.ResponseWriter, r *http.Request, user database.User) {
	name := r.URL.Query().Get("name")
	artist := r.URL.Query().Get("artist")

	if name != "" && artist != "" {
		music, err := apiCfg.DB.GetMusicByTitleAndArtist(r.Context(), database.GetMusicByTitleAndArtistParams{
			Name:   name,
			Artist: artist,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error retriving music by the artist: %v and song title: %v.\n Error: %s", artist, name, err))
		}
		respondWithJSON(w, 201, databaseMusicToMusic(music))
	} else if name != "" {
		music, err := apiCfg.DB.GetMusicByTitle(r.Context(), name)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error retriving music by the song title: %v.\n Error: %s", name, err))
		}
		respondWithJSON(w, 201, databaseMusicsToMusics(music))
	} else if artist != "" {
		music, err := apiCfg.DB.GetMusicByArtist(r.Context(), artist)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error retriving music by the artist: %v.\n Error: %s", artist, err))
		}
		respondWithJSON(w, 201, databaseMusicsToMusics(music))
	} else {
		respondWithError(w, 400, fmt.Sprintln("Error retriving music. Either song title, artist, or both are missing."))
	}
}

// TODO - get today's music
// 		  add caching mechanism to first check today's songs in the cache
// 		  before fetching from the database
// 		  Get today's saved data which gets saved everyday
