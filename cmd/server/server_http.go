package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
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
	router := chi.NewRouter()

	routing.SetupRoutes(router)

	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json",
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
	router.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	return &HttpServer{
		logger: logger,
		cfg:    cfg,
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
	}
}
