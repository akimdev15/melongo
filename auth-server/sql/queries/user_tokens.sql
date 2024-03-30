-- name: CreateUserToken :one
INSERT INTO user_tokens (id, api_key, access_token, refresh_token, expire_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserTokenByAPIKey :one
SELECT * FROM user_tokens WHERE api_key = $1;

-- name: UpdateToken :exec
UPDATE user_tokens
SET access_token = $1,
    refresh_token = $2,
    expire_time = $3,
    updated_at = $4
WHERE id = $5;
