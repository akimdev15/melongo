// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user_tokens.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUserToken = `-- name: CreateUserToken :one
INSERT INTO user_tokens (id, api_key, access_token, refresh_token, expire_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, api_key, access_token, refresh_token, expire_time, created_at, updated_at
`

type CreateUserTokenParams struct {
	ID           uuid.UUID
	ApiKey       string
	AccessToken  string
	RefreshToken string
	ExpireTime   time.Time
}

func (q *Queries) CreateUserToken(ctx context.Context, arg CreateUserTokenParams) (UserToken, error) {
	row := q.db.QueryRowContext(ctx, createUserToken,
		arg.ID,
		arg.ApiKey,
		arg.AccessToken,
		arg.RefreshToken,
		arg.ExpireTime,
	)
	var i UserToken
	err := row.Scan(
		&i.ID,
		&i.ApiKey,
		&i.AccessToken,
		&i.RefreshToken,
		&i.ExpireTime,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserTokenByAPIKey = `-- name: GetUserTokenByAPIKey :one
SELECT id, api_key, access_token, refresh_token, expire_time, created_at, updated_at FROM user_tokens WHERE api_key = $1
`

func (q *Queries) GetUserTokenByAPIKey(ctx context.Context, apiKey string) (UserToken, error) {
	row := q.db.QueryRowContext(ctx, getUserTokenByAPIKey, apiKey)
	var i UserToken
	err := row.Scan(
		&i.ID,
		&i.ApiKey,
		&i.AccessToken,
		&i.RefreshToken,
		&i.ExpireTime,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}