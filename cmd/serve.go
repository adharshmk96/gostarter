/*
Copyright © 2024 Adharsh M
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gostarter/cmd/server"
	"gostarter/infra/config"
	"gostarter/infra/logging"
	"log"
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

		logger, err := logging.NewFileLogger("logs/app.log")
		if err != nil {
			log.Fatal("failed to initialize logger...", err)
			return
		}

		svr := server.NewHttpServer(cfg, logger)

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-stop
			logger.Info("shutting down server...")
			err := svr.Stop(context.Background())
			if err != nil {
				logger.Error("failed to stop server", err)
				return
			}
		}()

		logger.Info("starting server...")
		fmt.Println("server starting at port:" + cfg.Server.Port)
		err = svr.Start()
		if err != nil {
			logger.Error("failed to start server", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

}
