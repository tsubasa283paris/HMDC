package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	DeckID          int32     `json:"deckId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	OwnerUserID     string    `json:"ownerUserId"`
	CurrentLeagueID NullInt32 `json:"currentLeagueId"`
	NumDuel         int32     `json:"numDuel"`
	NumWin          int32     `json:"numWin"`
}

type DeckStatsPerLeague struct {
	LeagueID int32 `json:"leagueId"`
	NumDuel  int32 `json:"numDuel"`
	NumWin   int32 `json:"numWin"`
}

type DeckDuelHistory struct {
	DuelID         int32     `json:"duelId"`
	LeagueID       NullInt32 `json:"leagueId"`
	OpponentUserID string    `json:"opponentUserId"`
	OpponentDeckID int32     `json:"opponentDeckId"`
	Result         string    `json:"result"`
	CreatedAt      time.Time `json:"createdAt"`
}

type DeckDetails struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	OwnerUserID     string    `json:"ownerUserId"`
	CurrentLeagueID NullInt32 `json:"currentLeagueId"`
}

type PutDeckDetailsParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

// Edit deck details: name, description and current league
func (c *Controller) PutDeckDetails(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("PutDeckDetails start")

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

	// receive body as API parameter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to read body",
			},
			errors.Wrap(err, "")
	}
	var param PutDeckDetailsParam
	err = json.Unmarshal(body, &param)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to decode body string to JSON format required by this API",
			},
			errors.Wrap(err, "")
	}
	log.Println("param:", string(body))

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

	// edit deck details
	err = queries.UpdateDeck(c.ctx, db.UpdateDeckParams{
		ID:          int32(paramDeckID),
		Name:        param.Name,
		Description: param.Description,
	})
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	responseBody := ErrorBody{
		Error: "",
	}

	log.Println("PutDeckDetails end")

	return http.StatusOK, responseBody, nil
}
