package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Return pointer to opened sql.DB
func DbCnx() (*sql.DB, error) {
	pgUser := os.Getenv("DB_USER")
	pgPassword := os.Getenv("DB_PASSWORD")
	if pgUser == "" || pgPassword == "" {
		return nil, errors.New("environment variable DB_USER or DB_PASSWORD not set")
	}

	pgTarget := fmt.Sprintf("user=%s password=%s dbname=hmdc sslmode=require", pgUser, pgPassword)
	dbCnx, err := sql.Open("postgres", pgTarget)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open connection with the database")
	}

	return dbCnx, nil
}
