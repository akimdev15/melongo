-- name: CreatePlaylist :one
INSERT INTO playlist (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetPlaylist :one
SELECT * 
FROM playlist 
WHERE user_id = $1 
AND name = $2;

-- name: GetPlaylistsByUserId :many
SELECT * 
FROM playlist
WHERE user_id = $1;
