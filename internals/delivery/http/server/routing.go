package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gostarter/infra"
	custommiddleware "gostarter/internals/delivery/http/middleware"
	"gostarter/internals/domain"
	"net/http"
	"strings"
)

func SetupRoutes(
	container *infra.Container,
	tokenService domain.TokenService,
	accountHandler domain.AccountHandler,
) *chi.Mux {
	cfg := container.Cfg

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Throttle(12000))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(custommiddleware.NewLatencyMiddleware(container.Meter))
	r.Use(custommiddleware.NewCounterMiddleware(container.Meter))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Routes
		accountRoutes(r, accountHandler, tokenService)
	})

	baseUrl := "http://" + strings.TrimPrefix(cfg.Server.BaseURL, "http://")

	// Swagger API docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: baseUrl + "/swagger/doc.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "gostarter API",
			},
			DarkMode: true,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(htmlContent))
	})
	r.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	return r

}

func accountRoutes(r chi.Router, accountHandler domain.AccountHandler, tokenService domain.TokenService) {
	r.Post("/auth/register", accountHandler.Register)
	r.Post("/auth/login", accountHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.JWTAuthMiddleware(tokenService))
		r.Post("/auth/logout", accountHandler.Logout)
		r.Get("/auth/profile", accountHandler.Profile)
	})
}
