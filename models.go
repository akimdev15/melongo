package main

import (
	"time"

	"github.com/akimdev15/melongo/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

type Music struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Artist    string    `json:"artist"`
	CreatedAt time.Time `json:"created_at"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

func databaseMusicToMusic(dbMusic database.Music) Music {
	return Music{
		ID:        dbMusic.ID,
		Name:      dbMusic.Name,
		Artist:    dbMusic.Artist,
		CreatedAt: dbMusic.CreatedAt,
	}
}
