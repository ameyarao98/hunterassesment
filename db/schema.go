package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func InitialiseSchema(db *pgx.Conn) error {
	if _, err := db.Exec(context.Background(),
		`
	CREATE TABLE IF NOT EXISTS "user"(
		username VARCHAR(100) PRIMARY KEY
	);
	`); err != nil {
		return err
	}
	return nil
}
