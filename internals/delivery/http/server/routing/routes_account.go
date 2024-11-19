package routing

import (
	custommiddleware "gostarter/internals/delivery/http/middleware"
	"gostarter/internals/delivery/http/web"
	"gostarter/internals/domain"

	"github.com/go-chi/chi/v5"
)

func accountApiRoutes(r chi.Router, accountHandler domain.AccountHandler) {
	r.Post("/auth/register", accountHandler.Register)
	r.Post("/auth/login", accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.IsAuthenticated)
		r.Post("/auth/logout", accountHandler.Logout)
		r.Get("/auth/profile", accountHandler.Profile)
	})
}

func accountWebRoutes(r chi.Router, handler *web.AccountWebHandler) {
	r.With(custommiddleware.RedirectIfLoggedIn("/profile")).Get("/register", handler.GetRegisterMember)
	r.Post("/register", handler.PostRegisterMember)
	r.With(custommiddleware.RedirectIfLoggedIn("/profile")).Get("/login", handler.GetLogin)
	r.Post("/login", handler.PostLogin)
	r.With(custommiddleware.IsAuthenticated).Get("/profile", handler.GetProfile)
	r.Post("/logout", handler.PostLogout)

}
