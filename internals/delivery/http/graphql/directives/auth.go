package directives

import (
	"context"
	"errors"
	"fitnesshub/internals/delivery/http/helpers"
	"fitnesshub/pkg/utils"
	"github.com/99designs/gqlgen/graphql"
	"slices"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	acc, err := helpers.GetAccountFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, errors.New("must be authenticated")
	}

	return next(ctx)
}

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*string) (interface{}, error) {
	acc, err := helpers.GetAccountFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, errors.New("must be authenticated")
	}

	rolesParsed := make([]string, len(roles))
	for i, r := range roles {
		rolesParsed[i] = utils.ParseNullString(r)
	}

	hasRole := slices.ContainsFunc(acc.Roles, func(r string) bool {
		return slices.Contains(rolesParsed, r)
	})
	if !hasRole {
		return nil, errors.New("unauthorized")
	}

	return next(ctx)
}
