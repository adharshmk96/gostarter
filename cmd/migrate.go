package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"gostarter/infra/config"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `This command allows you to run or revert database migrations.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("down")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
}

func runMigration(direction string) {
	cfg := config.NewConfig().Database
	dbURL := cfg.Postgres.Connection
	migrationsPath := "file://" + cfg.MigrationFiles

	db := postgres.Postgres{}
	drv, err := db.Open(dbURL)
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsPath, "pgx", drv)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	if direction == "up" {
		if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration up completed successfully")
	} else if direction == "down" {
		if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration down completed successfully")
	}
}
