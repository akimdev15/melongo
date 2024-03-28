-- +goose Up

CREATE TABLE user_tokens (
    id TEXT PRIMARY KEY,
    api_key VARCHAR(64) NOT NULL REFERENCES users(api_key),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expire_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX api_key_index ON user_tokens(api_key);

-- +goose Down

DROP TABLE user_tokens;

DROP INDEX IF EXISTS api_key_index;
