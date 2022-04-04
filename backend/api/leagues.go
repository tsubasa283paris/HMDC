package api

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"
)

// Get list of all leagues, containing only id and name
func (c *Controller) GetLeagues(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetLeagues start")

	// open database connection
	dbCnx, err := utils.DbCnx()
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to connect to the database",
			},
			errors.Wrap(err, "")
	}

	// prepare for query
	queries := db.New(dbCnx)

	// run query
	leagueList, err := queries.ListLeagues(c.ctx)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	log.Println("GetLeagues end")

	return http.StatusOK, leagueList, nil
}
