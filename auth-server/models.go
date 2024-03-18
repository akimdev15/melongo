package main

import (
	"time"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	APIKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		UserID:    dbUser.UserID,
		Email:     dbUser.Email,
		APIKey:    dbUser.ApiKey,
	}
}
