/*
Copyright © 2024 Adharsh Manikandan <debugslayer@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"gostarter/infra"
	"gostarter/infra/config"
	"gostarter/infra/pgdatabase"
	"gostarter/internals/domain"
	"gostarter/internals/service"
	"gostarter/internals/storage/pgstorage"
	"gostarter/pkg/testUtils"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resource",
}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create a new admin account",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		sqlConn := pgdatabase.NewConnection(cfg.Database.Postgres.Connection)
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		tracer := testUtils.NewNoopTracer()

		container := &infra.Container{
			Cfg:    cfg,
			Logger: logger,
			DbConn: sqlConn,
			Tracer: tracer,
		}

		accountRepo := pgstorage.NewAccountRepository(container)
		accountService := service.NewAccountService(container, accountRepo)

		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		acc := &domain.Account{
			Username: email,
			Email:    email,
			Password: password,
			Roles:    []string{domain.ROLE_ADMIN},
		}

		accountService.Register(context.Background(), acc)

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(adminCmd)

	adminCmd.Flags().StringP("email", "e", "", "Email of the admin")
	adminCmd.Flags().StringP("password", "p", "", "Password of the admin")

}
