package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ameyarao98/hunterassesment/factory/graph/generated"
	"github.com/ameyarao98/hunterassesment/factory/graph/model"
	pgx "github.com/jackc/pgx/v4"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context) (bool, error) {
	rows, err := r.DB.Query(context.Background(), `
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
		userResourceRows = append(userResourceRows, []any{ctx.Value(UserCtxKey), resources[i]})
	}

	copyCount, err := r.DB.CopyFrom(
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

// GetResourceData is the resolver for the getResourceData field.
func (r *queryResolver) GetResourceData(ctx context.Context) ([]*model.Resource, error) {
	rows, err := r.DB.Query(context.Background(), `
	SELECT resource_name, probability, color FROM "resource"`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		var resource model.Resource

		if err := rows.Scan(&resource.Name, &resource.Probability, &resource.Color); err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return resources, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
