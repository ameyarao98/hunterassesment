package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ameyarao98/hunterassesment/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

var resources = [3]string{"iron", "copper", "gold"}
var conn *pgx.Conn

func main() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	if err := db.InitialiseSchema(conn); err != nil {
		panic(err)
	}
	fmt.Println("Initilised db schema")

	go update()

	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yay"))
	})

	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		var userInput db.User
		err := json.NewDecoder(r.Body).Decode(&userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := db.CreateUser(conn, userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, resource := range resources {
			_, err := db.CreateUserResource(conn, db.UserResource{ResourceName: resource, Username: userInput.Username})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		var userInput db.User
		err := json.NewDecoder(r.Body).Decode(&userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := db.GetDashboardData(conn, userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)

	})

	fmt.Println("Server running on localhost:8080")
	http.ListenAndServe(":8080", r)
}

func update() {
	for range time.Tick(time.Second * 1) {
		go func() {
			err := db.UpdateResources(conn)
			if err != nil {
				fmt.Println(fmt.Errorf("update failed : %w", err))
				return
			}
			fmt.Println("update succeeded :)")

		}()
	}
}
