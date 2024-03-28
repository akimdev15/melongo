// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"
)

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	ApiKey    string
}

type UserToken struct {
	ID           string
	ApiKey       string
	AccessToken  string
	RefreshToken string
	ExpireTime   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
