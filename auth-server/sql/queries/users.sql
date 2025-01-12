-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, api_key)
VALUES ($1, $2, $3, $4, $5, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;
