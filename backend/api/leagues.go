package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"
)

type LeagueDuelHistory struct {
	DuelID    int32     `json:"duel_id"`
	LeagueID  NullInt32 `json:"league_id"`
	User1ID   string    `json:"user_1_id"`
	User2ID   string    `json:"user_2_id"`
	Deck1ID   int32     `json:"deck_1_id"`
	Deck2ID   int32     `json:"deck_2_id"`
	Result    int32     `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

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

// Get recent duel history with limit
func (c *Controller) GetLeagueDuelHistory(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetLeagueDuelHistory start")

	var err error

	// acquire query parameter
	paramLimitStr := r.FormValue("limit")
	var paramLimit int
	if paramLimitStr == "" {
		paramLimit = 100 // default value
	} else {
		// cast check
		paramLimit, err = strconv.Atoi(paramLimitStr)
		if err != nil {
			return http.StatusBadRequest,
				ErrorBody{
					Error: "failed to parse given query parameter to integer: limit",
				},
				errors.Wrap(err, "")
		}
	}
	log.Println("param limit:", paramLimit)

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
	duelHistoryList, err := queries.ListDuelsWithLimit(c.ctx, int32(paramLimit))
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []LeagueDuelHistory
	for _, duelHistory := range duelHistoryList {
		responseBody = append(responseBody, LeagueDuelHistory{
			DuelID:    duelHistory.ID,
			LeagueID:  NullInt32(duelHistory.LeagueID),
			User1ID:   duelHistory.User1ID,
			User2ID:   duelHistory.User2ID,
			Deck1ID:   duelHistory.Deck1ID,
			Deck2ID:   duelHistory.Deck2ID,
			Result:    duelHistory.Result,
			CreatedAt: duelHistory.CreatedAt,
		})
	}

	log.Println("GetLeagueDuelHistory end")

	return http.StatusOK, responseBody, nil
}
