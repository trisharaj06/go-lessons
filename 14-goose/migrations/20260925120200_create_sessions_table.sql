-- +goose Up
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash CHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_sessions_refresh_token_hash ON sessions (refresh_token_hash);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);

-- +goose Down
DROP TABLE IF EXISTS sessions;
