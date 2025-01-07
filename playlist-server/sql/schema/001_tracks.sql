-- +goose Up
CREATE TABLE tracks (
    rank INTEGER NOT NULL,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    uri TEXT NOT NULL,
    date DATE NOT NULL,
    PRIMARY KEY (uri, date)
);
CREATE INDEX idx_tracks_date ON tracks(date);

CREATE TABLE missed_tracks (
    rank INTEGER NOT NULL,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    date DATE NOT NULL,
	PRIMARY KEY (title, artist)
);

CREATE TABLE resolved_tracks (
	missed_title TEXT NOT NULL,
	missed_artist TEXT NOT NULL,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    uri TEXT NOT NULL,
	date DATE NOT NULL,
	PRIMARY KEY (missed_title, missed_artist)
);

-- +goose Down

DROP TABLE tracks;
DROP TABLE missed_tracks;
DROP TABLE resolved_tracks;
