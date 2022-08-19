package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ameyarao98/hunterassesment/factory/pb"
	"github.com/jackc/pgx/v4/pgxpool"
	grpc "google.golang.org/grpc"
)

var conn *pgxpool.Pool

type server struct {
	pb.UnimplementedFactoryServer
}

func (s *server) GetFactoryData(ctx context.Context, in *pb.GetFactoryDataRequest) (*pb.GetFactoryDataResponse, error) {
	factoryDatas, err := getFactoryData()
	if err != nil {
		return nil, err
	}
	return &pb.GetFactoryDataResponse{
		FactoryDatas: factoryDatas,
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
		`CREATE TABLE IF NOT EXISTS "resource"(resource_name VARCHAR(6) PRIMARY KEY);
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

func getFactoryData() ([]*pb.FactoryData, error) {
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
