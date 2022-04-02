package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

// Get list of all users, containing only id and name
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUsers start")

	// open database connection
	dbCnx, err := utils.DbCnx()
	if err != nil {
		RespondError(
			w,
			"failed to connect to the database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// prepare for query
	queries := db.New(dbCnx)

	// run query
	userList, err := queries.ListUsers(h.ctx)
	if err != nil {
		RespondError(
			w,
			"failed to communicate with database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// write response
	respondJSON(w, userList)

	log.Println("GetUsers end")
}
