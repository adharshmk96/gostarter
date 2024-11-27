package server

import (
	_ "gostarter/docs"
	"gostarter/infra"
	"gostarter/internals/delivery/http/graphql"
	"gostarter/internals/delivery/http/server/routing"
	"gostarter/internals/di"

	"context"
	"net/http"
)

type HttpServer struct {
	server *http.Server
}

func (s *HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func NewHttpServer(container *infra.Container) *HttpServer {
	storageDi := di.NewRepoContainer(container)
	serviceDi := di.NewServiceContainer(container, storageDi)
	handlerDi := di.NewHandlerContainer(container, serviceDi)

	r := routing.SetupRoutes(
		container,
		serviceDi,
		handlerDi,
	)

	gqlHandler := graphql.NewGQLHandler(container, serviceDi)
	gqlHandler.SetupRoutes(r)

	return &HttpServer{
		server: &http.Server{
			Addr:    ":" + container.Cfg.Server.Port,
			Handler: r,
		},
	}
}
