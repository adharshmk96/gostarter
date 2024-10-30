package routing

import (
	"github.com/go-chi/chi/v5"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/delivery/http/middleware"
	"gostarter/internals/service"
	"gostarter/internals/storage/memory"
)

func SetupRoutes(cfg *config.Config) func(r chi.Router) {
	return func(r chi.Router) {
		tokenService := service.NewTokenService(cfg.JWT)

		accountRepo := memory.NewAccountRepository()
		accountService := service.NewAccountService(accountRepo)
		accountHandler := api.NewAccountHandler(accountService, tokenService)

		r.Post("/v1/auth/register", accountHandler.Register)
		r.Post("/v1/auth/login", accountHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware(tokenService))
			r.Post("/v1/auth/logout", accountHandler.Logout)
			r.Get("/v1/auth/profile", accountHandler.Profile)
		})
	}
}
