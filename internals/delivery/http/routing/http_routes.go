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
		tokenService, err := service.NewTokenService(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath, cfg.JWT.ExpirationHours)
		if err != nil {
			panic(err)
		}

		// Storage
		accountRepo := memory.NewAccountRepository(container)

		// Services
		accountService := service.NewAccountService(container, accountRepo)

		// Handlers
		accountHandler := api.NewAccountHandler(container, accountService, tokenService)

		// Routes
		accountRoutes(r, accountHandler, tokenService)
	}
}
