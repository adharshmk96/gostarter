package routing

import (
	"github.com/go-chi/chi/v5"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/delivery/http/middleware"
	"gostarter/internals/domain"
)

func accountRoutes(r chi.Router, accountHandler *api.AccountHandler, tokenService domain.TokenService) {
	r.Post("/auth/register", accountHandler.Register)
	r.Post("/auth/login", accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(tokenService))
		r.Post("/auth/logout", accountHandler.Logout)
		r.Get("/auth/profile", accountHandler.Profile)
	})
}
