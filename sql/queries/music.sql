-- name: CreateMusic :one
INSERT INTO music (id, name, artist, created_at, genre)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMusicByTitle :many
SELECT * FROM music WHERE name = $1; 

-- name: GetMusicByArtist :many
SELECT * FROM music WHERE artist = $1;

-- name: GetMusicByTitleAndArtist :one
SELECT * FROM music WHERE name = $1 and artist = $2;