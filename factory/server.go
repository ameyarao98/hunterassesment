package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/ameyarao98/hunterassesment/factory/graph"
	"github.com/ameyarao98/hunterassesment/factory/graph/generated"
	"github.com/ameyarao98/hunterassesment/factory/graph/model"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func main() {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err := initiliaseSchema(conn); err != nil {
		panic(err)
	}
	log.Println("Schema initialised")
	router := chi.NewRouter()

	router.Use(func() func(http.Handler) http.Handler {
		authPublicKey, err := jwk.ParseKey([]byte(os.Getenv("AUTH_PUBLIC_KEY")))
		if err != nil {
			panic(err)
		}
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				verifiedToken, err := jwt.Parse([]byte(r.Header.Get("Creds")), jwt.WithKey(jwa.RS256, authPublicKey))
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), graph.UserCtxKey, verifiedToken.PrivateClaims()["id"])
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			})
		}
	}())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: conn}}))
	router.Handle("/graphql", srv)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}

}

func initiliaseSchema(conn *pgxpool.Pool) error {
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

	var resources []model.Resource
	if err := json.Unmarshal(configJsonBytes, &resources); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, resource := range resources {
		batch.Queue("INSERT INTO resource(resource_name, probability, color) VALUES($1, $2, $3) ON CONFLICT (resource_name) DO NOTHING;", resource.Name, resource.Probability, resource.Color)
	}
	br := conn.SendBatch(context.Background(), batch)
	if err := br.Close(); err != nil {
		return err
	}

	return nil
}
