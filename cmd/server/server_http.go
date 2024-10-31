package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "gostarter/docs"
	"gostarter/infra"
	custommiddleware "gostarter/internals/delivery/http/middleware"
	"gostarter/internals/delivery/http/routing"

	"context"
	"net/http"
)

type HttpServer struct {
	server *http.Server
}

func (s *HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func NewHttpServer(container *infra.Container) *HttpServer {
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

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// API Routes
	r.Route("/api/v1", routing.SetupRoutes(container))

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
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}
}
