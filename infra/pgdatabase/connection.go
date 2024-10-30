package pgdatabase

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(connStr string) (*sql.DB, error) {
	// Open the connection using pgx driver
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
