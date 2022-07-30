package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func InitialiseSchema(db *pgx.Conn) error {
	if _, err := db.Exec(context.Background(),
		`
		CREATE TABLE IF NOT EXISTS "user"(username VARCHAR(100) PRIMARY KEY);

		CREATE TABLE IF NOT EXISTS "resource"(resource_name VARCHAR(6) PRIMARY KEY);
		
		INSERT INTO
			"resource" (resource_name)
		VALUES
			('iron'),
			('copper'),
			('gold') ON CONFLICT (resource_name) DO NOTHING;
		
		CREATE TABLE IF NOT EXISTS "factory"(
			id SERIAL PRIMARY KEY,
			resource_name VARCHAR(6) NOT NULL REFERENCES "resource" ON DELETE CASCADE,
			factory_level INTEGER NOT NULL,
			production_per_second INTEGER NOT NULL,
			next_upgrade_duration INTEGER NOT NULL
		);
		
		INSERT INTO
			"factory" (
				id,
				resource_name,
				factory_level,
				production_per_second,
				next_upgrade_duration
			)
		VALUES
			(1, 'iron', 1, 10, 15),
			(2, 'iron', 2, 20, 30),
			(3, 'iron', 3, 40, 60),
			(4, 'iron', 4, 80, 90),
			(5, 'iron', 5, 150, 120),
			(6, 'copper', 1, 3, 15),
			(7, 'copper', 2, 7, 30),
			(8, 'copper', 3, 14, 60),
			(9, 'copper', 4, 30, 90),
			(10, 'copper', 5, 60, 120),
			(11, 'gold', 1, 2, 15),
			(12, 'gold', 2, 3, 30),
			(13, 'gold', 3, 4, 60),
			(14, 'gold', 4, 6, 90),
			(15, 'gold', 5, 8, 120) ON CONFLICT (id) DO NOTHING;

			CREATE TABLE IF NOT EXISTS "user_resource"(
				resource_name VARCHAR(6) NOT NULL REFERENCES "resource" ON DELETE CASCADE,
				username VARCHAR(100) NOT NULL REFERENCES "user" ON DELETE CASCADE,
				factory_level INTEGER DEFAULT 1 NOT NULL,
				amount INTEGER DEFAULT 0 NOT NULL,
				time_until_upgrade_complete INTEGER,
				PRIMARY KEY (resource_name,username)
			);
	`); err != nil {
		return err
	}
	return nil
}
