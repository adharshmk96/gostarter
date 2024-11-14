package routing

import (
	"github.com/go-chi/chi/v5"
	custommiddleware "gostarter/internals/delivery/http/middleware"
	"gostarter/internals/delivery/http/web"
	"gostarter/internals/domain"
)

func accountApiRoutes(r chi.Router, accountHandler domain.AccountHandler, tokenService domain.TokenService) {
	r.Post("/auth/register", accountHandler.Register)
	r.Post("/auth/login", accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.JWTAuthMiddleware(tokenService))
		r.Post("/auth/logout", accountHandler.Logout)
		r.Get("/auth/profile", accountHandler.Profile)
	})
}

func accountWebRoutes(r chi.Router, handler *web.AccountWebHandler, tokenService domain.TokenService) {
	r.With(custommiddleware.RedirectIfLoggedIn(tokenService, "/profile")).Get("/register", handler.GetRegisterMember)
	r.Post("/register", handler.PostRegisterMember)
	r.With(custommiddleware.RedirectIfLoggedIn(tokenService, "/profile")).Get("/login", handler.GetLogin)
	r.Post("/login", handler.PostLogin)
	r.With(custommiddleware.JWTAuthMiddleware(tokenService)).Get("/profile", handler.GetProfile)
	r.Post("/logout", handler.PostLogout)

}
