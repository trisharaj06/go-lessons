package models

import (
	"time"

	"goose-neon-demo/db"
)

type MagicLinkToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func CreateMagicLinkToken(userID int64, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO magic_link_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, userID, tokenHash, expiresAt)
	return err
}

// FindValidMagicLinkTokenByHash is read-only and must never mutate state -
// it's used by GET /auth/verify, which email scanners can prefetch.
func FindValidMagicLinkTokenByHash(tokenHash string) (*MagicLinkToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM magic_link_tokens
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()
	`

	var t MagicLinkToken
	err := db.DB.QueryRow(query, tokenHash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ConsumeMagicLinkToken atomically validates and marks a token used in one
// round trip, so two concurrent /auth/exchange calls can't both succeed.
func ConsumeMagicLinkToken(tokenHash string) (*MagicLinkToken, error) {
	query := `
		UPDATE magic_link_tokens
		SET used_at = NOW()
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()
		RETURNING id, user_id, token_hash, expires_at, used_at, created_at
	`

	var t MagicLinkToken
	err := db.DB.QueryRow(query, tokenHash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
