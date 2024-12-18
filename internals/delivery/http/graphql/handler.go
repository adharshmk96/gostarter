package graphql

import (
	"gostarter/infra"
	"gostarter/internals/delivery/http/graphql/directives"
	"gostarter/internals/delivery/http/graphql/generated"
	"gostarter/internals/delivery/http/graphql/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

type GQLHandler struct {
	config generated.Config
}

func NewGQLHandler(
	container *infra.Container,
	serviceDi *di.ServiceContainer,
) *GQLHandler {
	config := generated.Config{
		Resolvers: &resolver.Resolver{
			Container: container,
			ServiceDi: serviceDi,
		},
		Directives: generated.DirectiveRoot{
			Auth:    directives.Auth,
			HasRole: directives.HasRole,
		},
	}
	return &GQLHandler{
		config: config,
	}
}

func (h *GQLHandler) SetupRoutes(r chi.Router) {

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(h.config))

	// Playground handler
	r.Get("/playground", playground.Handler("Fitness Hub Graphql Server", "/query"))
	r.Post("/query", srv.ServeHTTP)
}
