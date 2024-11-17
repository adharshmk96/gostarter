package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"gostarter/infra"
	"gostarter/internals/delivery/http/graphql/generated"
	"gostarter/internals/delivery/http/graphql/resolver"
	"gostarter/internals/domain"
)

type GQLHandler struct {
	AccountService domain.AccountService
}

func NewGQLHandler(
	container *infra.Container,
	accountService domain.AccountService,
) *GQLHandler {
	return &GQLHandler{
		AccountService: accountService,
	}
}

func (h *GQLHandler) SetupRoutes(r chi.Router) {
	config := generated.Config{Resolvers: &resolver.Resolver{
		AccountService: h.AccountService,
	}}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	// Playground handler
	r.Get("/playground", playground.Handler("Fitness Hub Graphql Server", "/query"))
	r.Post("/query", srv.ServeHTTP)
}
