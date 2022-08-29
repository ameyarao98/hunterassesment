package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ameyarao98/hunterassesment/factory/pb"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	grpc "google.golang.org/grpc"
)

var conn *pgxpool.Pool

type server struct {
	pb.UnimplementedFactoryServer
}

func (s *server) GetFactoryData(ctx context.Context, in *pb.GetFactoryDataRequest) (*pb.GetFactoryDataResponse, error) {
	factoryDatas, err := getFactoryData(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetFactoryDataResponse{
		FactoryDatas: factoryDatas,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	created, err := createUser(ctx, int(in.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		Created: created,
	}, nil
}

func (s *server) UpgradeFactory(ctx context.Context, in *pb.UpgradeFactoryRequest) (*pb.UpgradeFactoryResponse, error) {
	upgraded, err := upgradeFactory(ctx, int(in.UserId), in.ResourceName)
	if err != nil {
		return nil, err
	}
	return &pb.UpgradeFactoryResponse{
		Upgraded: upgraded,
	}, nil
}
func (s *server) GetUserResourceData(ctx context.Context, in *pb.GetUserResourceDataRequest) (*pb.GetUserResourceDataResponse, error) {
	userResourceDatas, err := getUserResourceData(ctx, int(in.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResourceDataResponse{
		UserResourceDatas: userResourceDatas,
	}, nil
}

func main() {
	var err error
	conn, err = pgxpool.Connect(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err := initiliaseSchema(); err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "factory")
	})
	s := grpc.NewServer()
	pb.RegisterFactoryServer(s, &server{})
	go func() {
		for range time.Tick(time.Second * 1) {
			conn.Exec(context.Background(),
				`UPDATE "user_resource"
				SET amount = amount + production_per_second, 
				time_until_upgrade_complete = CASE time_until_upgrade_complete WHEN 1 THEN NULL ELSE "user_resource".time_until_upgrade_complete - 1 END, 
				factory_level = CASE time_until_upgrade_complete WHEN 1 THEN "user_resource".factory_level + 1 ELSE "user_resource".factory_level END
				FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level`)
		}
	}()
	fmt.Println("Server running")
	if err := s.Serve(listener); err != nil {
		panic(err)
	}

}

func initiliaseSchema() error {
	if _, err := conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS "resource"(resource_name VARCHAR(20) PRIMARY KEY);
		INSERT INTO "resource" (resource_name) VALUES ('iron'),('copper'),('gold') ON CONFLICT (resource_name) DO NOTHING;`,
	); err != nil {
		return err
	}

	if _, err := conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS "user_resource"(
			user_id INTEGER,
			resource_name VARCHAR(50) REFERENCES "resource" ON DELETE CASCADE,
			factory_level INTEGER DEFAULT 1 NOT NULL CHECK (factory_level <= 5),
			amount INTEGER DEFAULT 0 NOT NULL CHECK (amount >= 0),
			time_until_upgrade_complete INTEGER,
			PRIMARY KEY (user_id, resource_name)
		);`,
	); err != nil {
		return err
	}
	if _, err := conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS "factory"(
			resource_name VARCHAR(6) REFERENCES "resource" ON DELETE CASCADE,
			factory_level INTEGER,
			production_per_second INTEGER NOT NULL,
			next_upgrade_duration INTEGER NOT NULL,
			upgrade_cost JSON,
			PRIMARY KEY (resource_name, factory_level)
		);
		
		INSERT INTO
			"factory" (
				resource_name,
				factory_level,
				production_per_second,
				next_upgrade_duration,
				upgrade_cost
			)
		VALUES
			(
				'iron',
				1,
				10,
				15,
				'{ "iron": 300, "copper": 100, "gold": 1 }'
			),
			(
				'iron',
				2,
				20,
				30,
				'{ "iron": 800, "copper": 250, "gold": 2 }'
			),
			(
				'iron',
				3,
				40,
				60,
				'{ "iron": 1600, "copper": 500, "gold": 4 }'
			),
			(
				'iron',
				4,
				80,
				90,
				'{ "iron": 3000, "copper": 1000, "gold": 8 }'
			),
			('iron', 5, 150, 120, '{}'),
			(
				'copper',
				1,
				3,
				15,
				'{ "iron": 200, "copper": 70, "gold": 0}'
			),
			(
				'copper',
				2,
				7,
				30,
				'{ "iron": 400, "copper": 150, "gold": 0}'
			),
			(
				'copper',
				3,
				14,
				60,
				'{ "iron": 800, "copper": 300, "gold": 0}'
			),
			(
				'copper',
				4,
				30,
				90,
				'{ "iron": 1600, "copper": 600, "gold": 0}'
			),
			('copper', 5, 60, 120, '{}'),
			(
				'gold',
				1,
				2,
				15,
				'{ "iron": 0, "copper": 100, "gold": 2}'
			),
			(
				'gold',
				2,
				3,
				30,
				'{ "iron": 0, "copper": 200, "gold": 4}'
			),
			(
				'gold',
				3,
				4,
				60,
				'{ "iron": 0, "copper": 400, "gold": 8}'
			),
			(
				'gold',
				4,
				6,
				90,
				'{ "iron": 0, "copper": 800, "gold": 16}'
			),
			('gold', 5, 8, 120, '{}') ON CONFLICT (resource_name, factory_level) DO NOTHING;`,
	); err != nil {
		return err
	}
	return nil
}

func getFactoryData(ctx context.Context) ([]*pb.FactoryData, error) {
	rows, err := conn.Query(context.Background(), `
	SELECT resource_name, factory_level, production_per_second, next_upgrade_duration
	FROM "factory"`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dataz []*pb.FactoryData
	for rows.Next() {
		var data pb.FactoryData

		if err := rows.Scan(&data.ResourceName, &data.FactoryLevel, &data.ProductionPerSecond, &data.NextUpgradeDuration); err != nil {
			return nil, err
		}
		dataz = append(dataz, &data)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dataz, nil
}

func createUser(ctx context.Context, userID int) (bool, error) {
	rows, err := conn.Query(context.Background(), `
	SELECT resource_name FROM "resource"`,
	)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	var resources []string
	for rows.Next() {
		var resource string

		if err := rows.Scan(&resource); err != nil {
			return false, err
		}
		resources = append(resources, resource)
	}

	if rows.Err() != nil {
		return false, err
	}

	var resourceRows [][]any

	for i := 0; i < len(resources); i++ {
		resourceRows = append(resourceRows, []any{resources[i], userID})
	}

	copyCount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"user_resource"},
		[]string{"resource_name", "user_id"},
		pgx.CopyFromRows(resourceRows),
	)
	if err != nil {
		return false, err
	}
	return copyCount > 0, nil

}

func upgradeFactory(ctx context.Context, userID int, resourceName string) (bool, error) {
	var factoryLevel int
	var timeUntilUpgradeComplete *int
	var upgradeCost map[string]int
	if err := conn.
		QueryRow(context.Background(),
			`SELECT "user_resource".factory_level, "user_resource".time_until_upgrade_complete, "factory".upgrade_cost
			FROM "factory" INNER JOIN "user_resource" ON "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
			WHERE "user_resource".user_id=$1 AND "user_resource".resource_name=$2`,
			userID, resourceName).
		Scan(&factoryLevel, &timeUntilUpgradeComplete, &upgradeCost); err != nil {
		return false, err
	}

	if timeUntilUpgradeComplete != nil {
		return false, errors.New("upgrade already in progress")
	}
	if factoryLevel == 5 {
		return false, errors.New("factory level cannot go higher")
	}
	_, err := conn.Exec(context.Background(),
		`UPDATE "user_resource"
		SET time_until_upgrade_complete = CASE "user_resource".resource_name WHEN $1 THEN "factory".next_upgrade_duration ELSE time_until_upgrade_complete END, 
		amount = amount - ($2::json ->> "factory".resource_name)::integer
		FROM "factory" WHERE "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
		AND "user_resource".user_id=$3`,
		resourceName,
		upgradeCost,
		userID,
	)
	if err != nil {
		return false, err
	}
	return true, err
}

func getUserResourceData(ctx context.Context, userID int) ([]*pb.UserResourceData, error) {
	rows, err := conn.Query(context.Background(), `
	SELECT "user_resource".resource_name,"user_resource".factory_level,"user_resource".amount,"factory".production_per_second,"user_resource".time_until_upgrade_complete
	FROM "factory" INNER JOIN "user_resource" ON "factory".resource_name="user_resource".resource_name AND "factory".factory_level="user_resource".factory_level
	WHERE "user_resource".user_id=$1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dataz []*pb.UserResourceData
	for rows.Next() {
		var data pb.UserResourceData

		if err := rows.Scan(&data.ResourceName, &data.FactoryLevel, &data.Amount, &data.ProductionRate, &data.TimeUntilUpgradeComplete); err != nil {
			return nil, err
		}
		dataz = append(dataz, &data)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dataz, nil
}
