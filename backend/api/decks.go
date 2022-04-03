package api

import (
	"log"
	"net/http"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type Deck struct {
	DeckID          int32     `json:"deck_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	OwnerUserID     string    `json:"owner_user_id"`
	CurrentLeagueID NullInt32 `json:"current_league_id"`
	NumDuel         int32     `json:"num_duel"`
	NumWin          int32     `json:"num_win"`
}

// Get list of all decks, containing name, owner, description, current league,
// number of duels and number of victory
func (c *Controller) GetDecks(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetDecks start")

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

	// get decks
	deckList, err := queries.ListDecks(c.ctx)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []Deck
	for _, deck := range deckList {
		responseBody = append(responseBody, Deck{
			DeckID:          deck.ID,
			Name:            deck.Name,
			Description:     deck.Description,
			OwnerUserID:     deck.UserID,
			CurrentLeagueID: NullInt32(deck.LeagueID),
			NumDuel:         deck.NumDuel,
			NumWin:          deck.NumWin,
		})
	}

	log.Println("GetDecks end")

	return http.StatusOK, responseBody, nil
}
