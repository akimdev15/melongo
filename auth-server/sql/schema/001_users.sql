-- +goose Up

CREATE TABLE users( 
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  user_id TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL
);

-- +goose Down

DROP TABLE users;
