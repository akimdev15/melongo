-- name: CreateTrack :one
INSERT INTO tracks (id, name, artist, spotify_id, created_at)
VALUES ($1, $2, $3, $4, $5)
	RETURNING *;

-- name: GetTrackByNameAndArtist :one
SELECT * FROM tracks WHERE name = $1 and artist = $2;
