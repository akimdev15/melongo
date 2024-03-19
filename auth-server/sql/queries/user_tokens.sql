-- name: CreateUserToken :one
INSERT INTO user_tokens (id, api_key, access_token, refresh_token, expire_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserTokenByAPIKey :one
SELECT * FROM user_tokens WHERE api_key = $1;
