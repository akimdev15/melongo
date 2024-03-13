-- +goose Up
CREATE TABLE playlist (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE playlist_songs (
	playlist_id INT REFERENCES playlist(id) ON DELETE CASCADE,
	music_id UUID REFERENCES music(id) ON DELETE CASCADE,
	PRIMARY KEY (playlist_id, music_id)
);

-- +goose Down
DROP TABLE playlist;

DROP TABLE playlist_songs;