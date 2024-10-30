package routing

import (
	"github.com/go-chi/chi/v5"
	"gostarter/infra"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/service"
	"gostarter/internals/storage/memory"
)

func SetupRoutes(container *infra.Container) func(r chi.Router) {
	cfg := container.Cfg
	return func(r chi.Router) {
		tokenService := service.NewTokenService(cfg.JWT)

		// Storage
		accountRepo := memory.NewAccountRepository(container)

		// Services
		accountService := service.NewAccountService(container, accountRepo)

		// Handlers
		accountHandler := api.NewAccountHandler(container, accountService, tokenService)

		accountRoutes(r, accountHandler, tokenService)
	}
}
