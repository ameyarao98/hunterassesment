package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ameyarao98/hunterassesment/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	if err := db.InitialiseSchema(conn); err != nil {
		panic(err)
	}
	fmt.Println("Initilised db schema")

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
		for _, resource := range [3]string{"iron", "copper", "gold"} {
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
