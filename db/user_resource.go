package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type UserResource struct {
	ResourceName             string `json:"resource_name"`
	Username                 string `json:"username"`
	FactoryLevel             int    `json:"factory_level"`
	Amount                   int    `json:"amount"`
	TimeUntilUpgradeComplete int    `json:"time_until_upgrade_complete"`
}

type UpgradeInput struct {
	ResourceName string `json:"resource_name"`
	Username     string `json:"username"`
}

func CreateUserResource(db *pgx.Conn, input UserResource) (*UserResource, error) {
	var userResource UserResource
	if err := db.
		QueryRow(context.Background(),
			`INSERT INTO "user_resource" (resource_name, username) VALUES ($1, $2)
			 RETURNING resource_name, username, factory_level, amount`,
			input.ResourceName, input.Username).
		Scan(&userResource.ResourceName, &userResource.Username, &userResource.FactoryLevel, &userResource.Amount); err != nil {
		return nil, err
	}
	return &userResource, nil
}

func UpdateResources(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(),
		`UPDATE "user_resource"
		SET amount = amount + production_per_second, time_until_upgrade_complete=time_until_upgrade_complete-1, factory_level = CASE time_until_upgrade_complete WHEN 1 THEN "user_resource".factory_level + 1 ELSE "user_resource".factory_level END
		FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
		`)
	return err
}

func Upgrade(db *pgx.Conn, input UpgradeInput) error {
	_, err := db.Exec(context.Background(),
		`UPDATE "user_resource"
		SET time_until_upgrade_complete = "factory".next_upgrade_duration
		FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
		AND "user_resource".username=$1 AND "user_resource".resource_name=$2`,
		input.Username,
		input.ResourceName,
	)
	return err
}
