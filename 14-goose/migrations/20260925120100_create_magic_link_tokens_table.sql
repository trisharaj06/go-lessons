-- +goose Up
CREATE TABLE magic_link_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash CHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_magic_link_tokens_token_hash ON magic_link_tokens (token_hash);
CREATE INDEX idx_magic_link_tokens_user_id ON magic_link_tokens (user_id);

-- +goose Down
DROP TABLE IF EXISTS magic_link_tokens;
