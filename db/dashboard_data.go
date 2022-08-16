package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DashboardData struct {
	ResourceName             string         `json:"resource_name"`
	Amount                   int            `json:"amount"`
	FactoryLevel             int            `json:"factory_level"`
	ProductionRate           int            `json:"production_rate"`
	UpgradeCost              map[string]int `json:"upgrade_cost"`
	TimeUntilUpgradeComplete *int           `json:"time_until_upgrade_complete"`
}

func GetDashboardData(db *pgxpool.Pool, input User) ([]DashboardData, error) {
	rows, err := db.Query(context.Background(), `
	SELECT "user_resource".resource_name,"user_resource".amount,"user_resource".factory_level,"factory".production_per_second,"factory".upgrade_cost, "user_resource".time_until_upgrade_complete
	FROM "factory" INNER JOIN "user_resource" ON "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
	WHERE "user_resource".username=$1`,
		input.Username,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dataz []DashboardData
	for rows.Next() {
		var data DashboardData

		if err := rows.Scan(&data.ResourceName, &data.Amount, &data.FactoryLevel, &data.ProductionRate, &data.UpgradeCost, &data.TimeUntilUpgradeComplete); err != nil {
			return nil, err
		}
		dataz = append(dataz, data)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dataz, nil
}
