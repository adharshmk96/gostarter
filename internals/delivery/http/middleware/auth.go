package middleware

import (
	"context"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/helpers"
	"gostarter/internals/domain"
	"net/http"
)

func JWTMiddleware(tokenService domain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			userJWT, err := r.Cookie(config.AUTH_COOKIE_NAME)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Send account to context
			account, err := tokenService.ExtractAccount(userJWT.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "account", account)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuth(r.Context()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAuth(ctx context.Context) bool {
	acc, err := helpers.GetAccountFromContext(ctx)
	return acc != nil && err == nil
}

func RedirectIfLoggedIn(redirect string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			if isAuth(r.Context()) {
				http.Redirect(w, r, redirect, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}
