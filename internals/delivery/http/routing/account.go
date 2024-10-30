package routing

import (
	"github.com/go-chi/chi/v5"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/delivery/http/middleware"
	"gostarter/internals/domain"
)

func accountRoutes(r chi.Router, accountHandler *api.AccountHandler, tokenService domain.TokenService) {
	r.Post("/v1/auth/register", accountHandler.Register)
	r.Post("/v1/auth/login", accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(tokenService))
		r.Post("/v1/auth/logout", accountHandler.Logout)
		r.Get("/v1/auth/profile", accountHandler.Profile)
	})
}
