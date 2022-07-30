package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type User struct {
	Username string `json:"username"`
}

func GetUserByUsername(db *pgx.Conn, username string) (*User, error) {
	var user User
	if err := db.
		QueryRow(context.Background(),
			`SELECT username
			 FROM "user"
			 WHERE username=$1`,
			username).
		Scan(&user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *pgx.Conn, input User) (*User, error) {
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
