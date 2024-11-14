package server

import (
	_ "gostarter/docs"
	"gostarter/infra"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/delivery/http/server/routing"
	"gostarter/internals/delivery/http/web"
	"gostarter/internals/service"
	"gostarter/internals/storage/pgstorage"
	"gostarter/pkg/rendering"

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

	renderer := rendering.NewHtmlRenderer(config.TEMPLATE_DIR)

	// Storage
	//accountRepo := memory.NewAccountRepository(container)
	accountDbRepo := pgstorage.NewAccountRepository(container)

	// Services
	accountService := service.NewAccountService(container, accountDbRepo)

	// API Handlers
	accountHandler := api.NewAccountHandler(container, accountService, tokenService)
	accountWebHandler := web.NewAccountWebHandler(renderer, tokenService, accountService)

	r := routing.SetupRoutes(
		container,
		tokenService,
		accountHandler,

		accountWebHandler,
	)

	return &HttpServer{
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}
}
