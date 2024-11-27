package routing

import (
	"gostarter/infra"
	"gostarter/infra/config"
	custommiddleware "gostarter/internals/delivery/http/middleware"
	"gostarter/internals/di"
	"net/http"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func SetupRoutes(
	container *infra.Container,
	serviceDi *di.ServiceContainer,
	handlerDi *di.HandlerContainer,
) *chi.Mux {
	cfg := container.Cfg

	r := chi.NewRouter()

	// Setup Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
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

	r.Use(custommiddleware.JWTMiddleware(serviceDi.TokenService))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(config.STATIC_DIR))))

	// Web Routes
	accountWebRoutes(r, handlerDi.AccountWebHandler)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Routes
		accountApiRoutes(r, handlerDi.AccountHandler)
	})

	baseUrl := "http://" + strings.TrimPrefix(cfg.Server.BaseURL, "http://")

	// Swagger API docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
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
