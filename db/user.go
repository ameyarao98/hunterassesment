package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Username string `json:"username"`
}

func CreateUser(db *pgxpool.Pool, input User) (*User, error) {
	var user User
	if err := db.
		QueryRow(context.Background(),
			`INSERT INTO "user" (username) VALUES ($1)
			 RETURNING username`,
			input.Username).
		Scan(&user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
