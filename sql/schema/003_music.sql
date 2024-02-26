-- +goose Up

CREATE TABLE music (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	artist TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

-- +goose Down

DROP TABLE music;