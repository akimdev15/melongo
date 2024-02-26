package main

import (
	"context"
	"fmt"
	"time"

	"github.com/akimdev15/melongo/internal/database"
	"github.com/akimdev15/mscraper"
	"github.com/google/uuid"
)

// handleDailyMusicScrape - save music to the database
func (apiCfg *apiConfig) handleDailyMusicScrape() {
	songs := mscraper.GetNewestHipHopFromMelon()

	for _, song := range songs {
		savedSong, err := apiCfg.DB.CreateMusic(context.Background(), database.CreateMusicParams{
			ID:        uuid.New(),
			Name:      song.Title,
			Artist:    song.Artist,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			fmt.Printf("Error creating the music: %v.\n  Error: %s", song, err)
		}
		fmt.Println("Saved song: ", savedSong)
	}

}

// TODO - get today's music
// 		  add caching mechanism to first check today's songs in the cache
// 		  before fetching from the database
// 		  Get today's saved data which gets saved everyday
