-- name: CreateMusic :one
INSERT INTO music (id, name, artist, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;