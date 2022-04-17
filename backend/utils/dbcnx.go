package utils

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
)

// Return pointer to opened sql.DB
func DbCnx() (*sql.DB, error) {
	pgUrl := os.Getenv("DATABASE_URL")
	if pgUrl == "" {
		return nil, errors.New("environment variable DATABASE_URL not set")
	}

	dbCnx, err := sql.Open("postgres", pgUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open connection with the database")
	}

	return dbCnx, nil
}
