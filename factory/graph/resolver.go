package graph

import "github.com/jackc/pgx/v4/pgxpool"

var UserCtxKey = &contextKey{"userID"}

type contextKey struct {
	name string
}

type Resolver struct {
	DB *pgxpool.Pool
}
