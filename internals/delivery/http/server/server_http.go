package server

import (
	_ "gostarter/docs"
	"gostarter/infra"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/service"
	"gostarter/internals/storage/memory"

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
	cfg := container.Cfg

	tokenService := service.NewTokenService(cfg.JWT)

	// Storage
	accountRepo := memory.NewAccountRepository(container)

	// Services
	accountService := service.NewAccountService(container, accountRepo)

	// Handlers
	accountHandler := api.NewAccountHandler(container, accountService, tokenService)

	r := SetupRoutes(
		container,
		tokenService,
		accountHandler,
	)

	return &HttpServer{
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}
}
