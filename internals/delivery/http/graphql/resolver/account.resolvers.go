package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"gostarter/internals/delivery/http/graphql/generated"
	"gostarter/internals/delivery/http/graphql/models"
	"gostarter/internals/domain"
)

// Roles is the resolver for the roles field.
func (r *accountResolver) Roles(ctx context.Context, obj *domain.Account) ([]models.Role, error) {
	var roles []models.Role
	for _, role := range obj.Roles {
		roles = append(roles, models.Role(role))
	}
	return roles, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *accountResolver) CreatedAt(ctx context.Context, obj *domain.Account) (string, error) {
	timeString := obj.CreatedAt.Format("2006-01-02 15:04:05")

	return timeString, nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *accountResolver) UpdatedAt(ctx context.Context, obj *domain.Account) (string, error) {
	timeString := obj.UpdatedAt.Format("2006-01-02 15:04:05")

	return timeString, nil
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

type accountResolver struct{ *Resolver }