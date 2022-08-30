package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/ameyarao98/hunterassesment/factory/pb"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	grpc "google.golang.org/grpc"
)

var conn *pgxpool.Pool

type Resource struct {
	ResourceName string  `json:"resource_name"`
	Probability  float32 `json:"probability"`
	Color        string  `json:"color"`
}
type server struct {
	pb.UnimplementedFactoryServer
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

func (s *server) GetResourceData(ctx context.Context, in *pb.GetResourceDataRequest) (*pb.GetResourceDataResponse, error) {
	resourceDataz, err := getResourceData(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetResourceDataResponse{
		ResourceDataz: resourceDataz,
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
	fmt.Println("Server running")
	if err := s.Serve(listener); err != nil {
		panic(err)
	}

}

func initiliaseSchema() error {
	if _, err := conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS "resource"(
			resource_name VARCHAR(20) PRIMARY KEY, 
			probability FLOAT(10), 
			color VARCHAR(6)
		  );
		  `,
	); err != nil {
		return err
	}

	if _, err := conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS "user_resource"(
			user_id INTEGER,
			resource_name VARCHAR(50) REFERENCES "resource" ON DELETE CASCADE,
			amount INTEGER DEFAULT 0 NOT NULL CHECK (amount >= 0),
			PRIMARY KEY (user_id, resource_name)
		);`,
	); err != nil {
		return err
	}

	configJson, err := os.Open("resources.json")
	if err != nil {
		return err
	}
	configJsonBytes, err := io.ReadAll(configJson)
	if err != nil {
		return err
	}

	defer configJson.Close()

	var resources []Resource
	if err := json.Unmarshal(configJsonBytes, &resources); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, resource := range resources {
		batch.Queue("INSERT INTO resource(resource_name, probability, color) VALUES($1, $2, $3) ON CONFLICT (resource_name) DO NOTHING;", resource.ResourceName, resource.Probability, resource.Color)
	}
	br := conn.SendBatch(context.Background(), batch)
	if err := br.Close(); err != nil {
		return err
	}

	return nil
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

	var userResourceRows [][]any

	for i := 0; i < len(resources); i++ {
		userResourceRows = append(userResourceRows, []any{userID, resources[i]})
	}

	copyCount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"user_resource"},
		[]string{"user_id", "resource_name"},
		pgx.CopyFromRows(userResourceRows),
	)
	if err != nil {
		return false, err
	}
	return copyCount > 0, nil

}
func getResourceData(ctx context.Context) ([]*pb.ResourceData, error) {
	rows, err := conn.Query(context.Background(), `
	SELECT resource_name, probability, color FROM "resource"`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dataz []*pb.ResourceData
	for rows.Next() {
		var data pb.ResourceData

		if err := rows.Scan(&data.ResourceName, &data.Probability, &data.Color); err != nil {
			return nil, err
		}
		dataz = append(dataz, &data)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dataz, nil
}
