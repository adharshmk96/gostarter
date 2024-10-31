package cmd

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"gostarter/infra/config"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		name, _ := cmd.Flags().GetString("name")
		ext, _ := cmd.Flags().GetString("ext")

		if err := createMigration(dir, name, ext); err != nil {
			log.Fatalf("Could not create migration: %v", err)
		}
		log.Println("Migration created successfully")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)

	migrateCreateCmd.Flags().StringP("name", "n", "", "Name of the migration")
	migrateCreateCmd.Flags().StringP("ext", "e", ".sql", "Extension of the migration file")
	migrateCreateCmd.Flags().StringP("dir", "d", "migrations", "Directory to store the migration file")
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

func createMigration(dir, name, ext string) error {
	dir = filepath.Clean(dir)
	ext = "." + strings.TrimPrefix(ext, ".")

	matches, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return err
	}

	nextSeqVersion, err := nextVersion(matches)
	if err != nil {
		return err
	}

	if hasDuplicates(nextSeqVersion, dir, ext) {
		return fmt.Errorf("duplicate migration nextSeqVersion: %s", nextSeqVersion)
	}

	// create the file
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", nextSeqVersion, name, direction, ext)
		filename := filepath.Join(dir, basename)
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
		log.Println("created migration file: ", filename)
	}

	return nil
}

func nextVersion(matches []string) (string, error) {
	if len(matches) == 0 {
		return "000001", nil
	}

	lastFile := matches[len(matches)-1]
	base := filepath.Base(lastFile)
	idx := strings.Index(base, "_")

	if idx < 1 {
		return "", fmt.Errorf("invalid migration file: %s", base)
	}

	seqString := base[:idx]
	seq, err := strconv.ParseUint(seqString, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", seq+1), nil
}

func hasDuplicates(version string, dir, ext string) bool {
	versionGlob := filepath.Join(dir, fmt.Sprintf("%s_*%s", version, ext))
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		return false
	}

	return len(matches) > 0
}
