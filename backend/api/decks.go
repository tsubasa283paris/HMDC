package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

type DeckDuelHistory struct {
	DuelID         int32     `json:"duel_id"`
	LeagueID       NullInt32 `json:"league_id"`
	OpponentUserID string    `json:"opponent_user_id"`
	OpponentDeckID int32     `json:"opponent_deck_id"`
	Result         string    `json:"result"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeckDetails struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	OwnerUserID     string    `json:"owner_user_id"`
	CurrentLeagueID NullInt32 `json:"current_league_id"`
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

// Get duel history of the specified user
func (c *Controller) GetDeckDuelHistory(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetDeckDuelHistory start")

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

	// get duel history
	duelHistoryList, err := queries.ListDeckDuelHistory(c.ctx, int32(paramDeckID))
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []DeckDuelHistory
	for _, duelHistory := range duelHistoryList {
		responseBody = append(responseBody, DeckDuelHistory{
			DuelID:         duelHistory.ID,
			LeagueID:       NullInt32(duelHistory.LeagueID),
			OpponentUserID: duelHistory.OpponentUserID,
			OpponentDeckID: duelHistory.OpponentDeckID,
			Result:         fmt.Sprintf("%v", duelHistory.Result),
			CreatedAt:      duelHistory.CreatedAt,
		})
	}

	log.Println("GetDeckDuelHistory end")

	return http.StatusOK, responseBody, nil
}

// Get deck details, containing name, owner, description and current league
func (c *Controller) GetDeckDetails(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetDeckDetails start")

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

	// get deck details
	deck, err := queries.GetDeck(c.ctx, int32(paramDeckID))
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

	// create responseBody
	responseBody := DeckDetails{
		Name:            deck.Name,
		Description:     deck.Description,
		OwnerUserID:     deck.UserID,
		CurrentLeagueID: NullInt32(deck.LeagueID),
	}

	log.Println("GetDeckDetails end")

	return http.StatusOK, responseBody, nil
}
