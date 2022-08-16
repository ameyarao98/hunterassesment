package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserResource struct {
	ResourceName             string `json:"resource_name"`
	Username                 string `json:"username"`
	FactoryLevel             int    `json:"factory_level"`
	Amount                   int    `json:"amount"`
	TimeUntilUpgradeComplete *int   `json:"time_until_upgrade_complete"`
}

type UpgradeInput struct {
	ResourceName string `json:"resource_name"`
	Username     string `json:"username"`
}

func CreateUserResource(db *pgxpool.Pool, input UserResource) (*UserResource, error) {
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

func UpdateResources(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`UPDATE "user_resource"
		SET amount = amount + production_per_second, 
		time_until_upgrade_complete = CASE time_until_upgrade_complete WHEN 1 THEN NULL ELSE "user_resource".time_until_upgrade_complete - 1 END, 
		factory_level = CASE time_until_upgrade_complete WHEN 1 THEN "user_resource".factory_level + 1 ELSE "user_resource".factory_level END
		FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
		`)
	return err
}

func Upgrade(db *pgxpool.Pool, input UpgradeInput) error {
	var userResource UserResource
	var upgradeCost map[string]int
	if err := db.
		QueryRow(context.Background(),
			`SELECT "user_resource".resource_name, "user_resource".username, "user_resource".factory_level,"user_resource".amount, "user_resource".time_until_upgrade_complete, "factory".upgrade_cost
			FROM "factory" INNER JOIN "user_resource" ON "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
			WHERE "user_resource".resource_name=$1 AND "user_resource".username=$2`,
			input.ResourceName, input.Username).
		Scan(&userResource.ResourceName, &userResource.Username, &userResource.FactoryLevel, &userResource.Amount, &userResource.TimeUntilUpgradeComplete, &upgradeCost); err != nil {
		return err
	}

	if userResource.TimeUntilUpgradeComplete != nil {
		return errors.New("upgrade already in progress")
	}
	if userResource.FactoryLevel == 5 {
		return errors.New("factory level cannot go higher")
	}
	_, err := db.Exec(context.Background(),
		`UPDATE "user_resource"
		SET time_until_upgrade_complete = CASE "user_resource".resource_name WHEN $1 THEN "factory".next_upgrade_duration ELSE time_until_upgrade_complete END, 
		amount = amount - ($2::json ->> "factory".resource_name)::integer
		FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
		AND "user_resource".username=$3`,
		input.ResourceName,
		upgradeCost,
		input.Username,
	)
	return err
}
