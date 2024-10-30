package middleware

import (
	"context"
	"gostarter/infra/config"
	"gostarter/internals/domain"
	"net/http"
)

func JWTAuthMiddleware(tokenService domain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			userJWT, err := r.Cookie(config.AUTH_COOKIE_NAME)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Send account to context
			account, err := tokenService.ExtractAccount(userJWT.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
