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

-- +goose Down

DROP TABLE tracks;
