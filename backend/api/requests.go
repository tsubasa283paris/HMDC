package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"
)

type DuelRequest struct {
	DuelID         int32     `json:"duel_id"`
	LeagueID       NullInt32 `json:"league_id"`
	OpponentUserID string    `json:"opponent_user_id"`
	DeckID         int32     `json:"deck_id"`
	OpponentDeckID int32     `json:"opponent_deck_id"`
	Result         string    `json:"result"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      int32     `json:"created_by"`
}

// Get unconfirmed duel history of the specified user
func (c *Controller) GetUserDuelRequests(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUserDuelRequests start")

	// acquire ID of the requested user
	reqUserID := r.Header.Get("UserID")

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

	// get requests (= duel histories that are not confirmed yet)
	UnconfirmedDuelList, err := queries.ListUnconfirmedDuelsByUser(c.ctx, reqUserID)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []DuelRequest
	for _, duelHistory := range UnconfirmedDuelList {
		responseBody = append(responseBody, DuelRequest{
			DuelID:         duelHistory.ID,
			LeagueID:       NullInt32(duelHistory.LeagueID),
			OpponentUserID: duelHistory.OpponentUserID,
			DeckID:         duelHistory.DeckID,
			OpponentDeckID: duelHistory.OpponentDeckID,
			Result:         fmt.Sprintf("%v", duelHistory.Result),
			CreatedAt:      duelHistory.CreatedAt,
			CreatedBy:      duelHistory.CreatedBy,
		})
	}

	log.Println("GetUserDuelRequests end")

	return http.StatusOK, responseBody, nil
}
