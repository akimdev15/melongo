-- +goose Up
CREATE TABLE tracks (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	artist TEXT NOT NULL,
	spotify_id TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	CONSTRAINT unique_track_name_artist UNIQUE (name, artist)
);

-- +goose Down

DROP TABLE tracks;
