-- +goose Up

CREATE TABLE genre (
    id int PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    code TEXT NOT NULL UNIQUE
);

CREATE TABLE music (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    artist TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    genre INT NOT NULL REFERENCES genre(id),
    CONSTRAINT unique_music_name_artist UNIQUE (name, artist)
);

INSERT INTO genre (id, name, code) 
VALUES (1, 'KBALLAD', '0100');
INSERT INTO genre (id, name, code) 
VALUES (2, 'KDANCE', '0200');
INSERT INTO genre (id, name, code) 
VALUES (3, 'KHIPHOP', '0300');
INSERT INTO genre (id, name, code) 
VALUES (4, 'KRB', '0400');
INSERT INTO genre (id, name, code) 
VALUES (5, 'KINDY', '0500');
INSERT INTO genre (id, name, code) 
VALUES (6, 'KROCK', '0600');
INSERT INTO genre (id, name, code) 
VALUES (7, 'KTROT', '0700');
INSERT INTO genre (id, name, code) 
VALUES (8, 'KBLUSE', '0800');
INSERT INTO genre (id, name, code) 
VALUES (9, 'POP', '0900');
INSERT INTO genre (id, name, code) 
VALUES (10, 'ROCK', '1000');
INSERT INTO genre (id, name, code) 
VALUES (11, 'ELEC', '1100');
INSERT INTO genre (id, name, code) 
VALUES (12, 'HIPHOP', '1200');
INSERT INTO genre (id, name, code) 
VALUES (13, 'RB', '1300');
INSERT INTO genre (id, name, code) 
VALUES (14, 'BLUSE', '1400');

-- +goose Down

DROP TABLE music;

DROP TABLE genre;