package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type UserStatsPerLeague struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

// Get list of all users, containing only id and name
func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUsers start")

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
	userList, err := queries.ListUsers(c.ctx)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	log.Println("GetUsers end")

	return http.StatusOK, userList, nil
}

// Get stats of the specified user
func (c *Controller) GetUserStats(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUserStats start")

	// acquire URL parameter
	paramUserID := chi.URLParam(r, "userId")
	if paramUserID == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "query parameter missing: userId",
			},
			nil
	}

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

	// check if the specified ID exists
	_, err = queries.GetUser(c.ctx, paramUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "user not found (blame: query parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// get stats
	statsList, err := queries.GetUserStats(c.ctx, paramUserID)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []UserStatsPerLeague
	for _, statsPerLeague := range statsList {
		responseBody = append(responseBody, UserStatsPerLeague{
			LeagueID: statsPerLeague.LeagueID,
			NumDuel:  statsPerLeague.NumDuel,
			NumWin:   statsPerLeague.NumWin,
		})
	}

	log.Println("GetUserStats end")

	return http.StatusOK, responseBody, nil
}
