package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

// Get list of all users, containing only id and name
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUsers start")

	// open database connection
	dbCnx, err := utils.DbCnx()
	if err != nil {
		errorResponse(
			w,
			"failed to connect to the database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// prepare for query
	ctx := context.Background()
	queries := db.New(dbCnx)

	// run query
	userList, err := queries.ListUsers(ctx)
	if err != nil {
		errorResponse(
			w,
			"failed to communicate with database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// generate body JSON
	bodyJSONBytes, err := json.Marshal(userList)
	if err != nil {
		errorResponse(
			w,
			"failed to convert data given by database to JSON",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// write response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bodyJSONBytes)

	log.Println("GetUsers end")
}
