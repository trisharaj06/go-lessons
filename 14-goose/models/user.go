package models

import (
	"time"

	"goose-neon-demo/db"
)

type User struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}

func FindOrCreateUserByEmail(email string) (*User, error) {
	query := `
		INSERT INTO users (email) VALUES ($1)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id, email, created_at
	`

	var u User
	err := db.DB.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func FindUserByID(id int64) (*User, error) {
	query := `SELECT id, email, created_at FROM users WHERE id = $1`

	var u User
	err := db.DB.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
