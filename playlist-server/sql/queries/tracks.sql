-- name: CreateTrack :one
INSERT INTO tracks (rank, title, artist, uri, date)
VALUES ($1, $2, $3, $4, $5)
	RETURNING *;

-- name: GetTracksByDate :many
SELECT * FROM tracks WHERE date = $1;


-- name: CreateMissedTrack :one
INSERT INTO missed_tracks (rank, title, artist, date)
VALUES ($1, $2, $3, $4)
	RETURNING *;

-- name: GetMissedTracks :one
SELECT * FROM missed_tracks WHERE title = $1 AND artist = $2;

-- name: GetMissedTracksByDate :many
SELECT * FROM missed_tracks WHERE date = $1;


-- name: CreateResolvedTrack :one
INSERT INTO resolved_tracks (missed_title, missed_artist, title, artist, uri, date)
VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *;

-- name: GetResolvedTrack :one
SELECT * FROM resolved_tracks WHERE missed_title = $1 AND missed_artist = $2;

-- name: RemoveMissedTrack :one
DELETE FROM missed_tracks WHERE title = $1 AND artist = $2
RETURNING *;