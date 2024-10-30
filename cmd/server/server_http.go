package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "gostarter/docs"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/routing"

	"context"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	logger *slog.Logger
	cfg    *config.Config
	server *http.Server
}

func (s *HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func NewHttpServer(cfg *config.Config, logger *slog.Logger) *HttpServer {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// API Routes
	r.Route("/api", routing.SetupRoutes(cfg))

	// Swagger API docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "http://" + cfg.Server.BaseURL + "/swagger/doc.json",
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

	return &HttpServer{
		logger: logger,
		cfg:    cfg,
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}
}
