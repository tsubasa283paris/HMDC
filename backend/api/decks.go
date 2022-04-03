package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

type DeckStatsPerLeague struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
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

// Get stats of the specified deck
func (c *Controller) GetDeckStats(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetDeckStats start")

	// acquire URL parameter
	paramDeckIDStr := chi.URLParam(r, "deckId")
	if paramDeckIDStr == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: deckId",
			},
			errors.New("url parameter missing")
	}
	paramDeckID, err := strconv.Atoi(paramDeckIDStr)
	if err != nil {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "failed to parse given url parameter to integer: deckId",
			},
			errors.Wrap(err, "")
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
	_, err = queries.GetDeck(c.ctx, int32(paramDeckID))
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "deck not found (blame: url parameter)",
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
	statsList, err := queries.ListDeckStats(c.ctx, int32(paramDeckID))
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []DeckStatsPerLeague
	for _, statsPerLeague := range statsList {
		responseBody = append(responseBody, DeckStatsPerLeague{
			LeagueID: statsPerLeague.LeagueID,
			NumDuel:  statsPerLeague.NumDuel,
			NumWin:   statsPerLeague.NumWin,
		})
	}

	log.Println("GetDeckStats end")

	return http.StatusOK, responseBody, nil
}
