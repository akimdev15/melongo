-- name: CreateTrack :one
INSERT INTO tracks (rank, title, artist, uri, date)
VALUES ($1, $2, $3, $4, $5)
	RETURNING *;

-- name: GetTracksByDate :many
SELECT * FROM tracks WHERE date = $1;
