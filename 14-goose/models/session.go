package models

import (
	"time"

	"goose-neon-demo/db"
)

type Session struct {
	ID               int64
	UserID           int64
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

func CreateSession(userID int64, refreshTokenHash string, expiresAt time.Time) (*Session, error) {
	query := `
		INSERT INTO sessions (user_id, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, refresh_token_hash, expires_at, created_at
	`

	var s Session
	err := db.DB.QueryRow(query, userID, refreshTokenHash, expiresAt).
		Scan(&s.ID, &s.UserID, &s.RefreshTokenHash, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func FindSessionByRefreshHash(refreshTokenHash string) (*Session, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, expires_at, created_at
		FROM sessions
		WHERE refresh_token_hash = $1 AND expires_at > NOW()
	`

	var s Session
	err := db.DB.QueryRow(query, refreshTokenHash).
		Scan(&s.ID, &s.UserID, &s.RefreshTokenHash, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSession(id int64) error {
	_, err := db.DB.Exec(`DELETE FROM sessions WHERE id = $1`, id)
	return err
}

// DeleteSessionByRefreshHash is a no-op (not an error) if no session matches,
// so /auth/logout is idempotent.
func DeleteSessionByRefreshHash(refreshTokenHash string) error {
	_, err := db.DB.Exec(`DELETE FROM sessions WHERE refresh_token_hash = $1`, refreshTokenHash)
	return err
}
