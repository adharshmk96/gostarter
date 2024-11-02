/*
Copyright © 2024 Adharsh M
*/
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gostarter/infra"
	"gostarter/infra/config"
	"gostarter/infra/logging"
	"gostarter/infra/observability"
	"gostarter/infra/pgdatabase"
	"gostarter/internals/delivery/http/server"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()

		logger := logging.NewFileLogger("logs/app.log")

		telemetry := observability.NewTelemetryService(&cfg.Observability)
		defer telemetry.Stop()
		// Tracers can be separated for modules in the future if needed
		tracer := telemetry.GetTracerProvider().Tracer(cfg.Observability.TracerName)
		meter := telemetry.GetMeterProvider().Meter(cfg.Observability.MeterName)

		sqlConn := pgdatabase.NewConnection(cfg.Database.Postgres.Connection)

		container := &infra.Container{
			Cfg:    cfg,
			DbConn: sqlConn,
			Logger: logger,
			Tracer: tracer,
			Meter:  meter,
		}

		svr := server.NewHttpServer(container)

		stop := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-stop
			logger.Info("shutting down server...")
			if err := svr.Stop(context.Background()); err != nil {
				logger.Error("failed to stop server:", slog.String("error", err.Error()))
			}
			done <- true
		}()

		logger.Info("starting server...")
		log.Printf("server starting at port: %s\n", cfg.Server.Port)

		if err := svr.Start(); err != nil {
			logger.Error("server stopped:", slog.String("error", err.Error()))
		}

		<-done

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

}
