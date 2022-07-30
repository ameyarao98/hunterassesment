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
