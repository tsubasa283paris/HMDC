package utils

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
)

// Return pointer to opened sql.DB
func DbCnx() (*sql.DB, error) {
	dbCnx, err := sql.Open(
		"postgres",
		"host="+os.Getenv("DB_HOST")+
			" port="+os.Getenv("DB_PORT")+
			" dbname="+os.Getenv("DB_NAME")+
			" user="+os.Getenv("DB_USER")+
			" password="+os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open connection with the database")
	}

	return dbCnx, nil
}
