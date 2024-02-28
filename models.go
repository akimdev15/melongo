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

type Genre struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
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

func databaseMusicsToMusics(dbMusics []database.Music) []Music {
	musics := []Music{}
	for _, dbMusic := range dbMusics {
		musics = append(musics, databaseMusicToMusic(dbMusic))
	}
	return musics
}

func databaseGenreToGenre(dbGenre database.Genre) Genre {
	return Genre{
		ID:   dbGenre.ID,
		Name: dbGenre.Name,
		Code: dbGenre.Code,
	}
}

func databaseGenresToGenres(dbGenres []database.Genre) []Genre {
	genres := []Genre{}
	for _, dbGenre := range dbGenres {
		genres = append(genres, databaseGenreToGenre(dbGenre))
	}
	return genres
}
