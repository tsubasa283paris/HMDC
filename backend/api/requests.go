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
	"github.com/pkg/errors"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"
)

// デュエル結果として想定するEnum
type DuelResult int

const (
	DuelResultDraw DuelResult = iota
	DuelResultWin
	DuelResultLose
	DuelResultNA
	DuelResultInvalid
)

func (dr DuelResult) String() string {
	switch dr {
	case DuelResultDraw:
		return "draw"
	case DuelResultWin:
		return "win"
	case DuelResultLose:
		return "lose"
	case DuelResultNA:
		return "n/a"
	case DuelResultInvalid:
		return "invalid"
	default:
		return "undefined"
	}
}

func (dr DuelResult) Int32() int32 {
	switch dr {
	case DuelResultDraw:
		return 0
	case DuelResultWin:
		return 1
	case DuelResultLose:
		return 2
	case DuelResultNA:
		return -1
	case DuelResultInvalid:
		return -1
	default:
		return -1
	}
}

func GetValidDuelResults() []DuelResult {
	return []DuelResult{
		DuelResultDraw,
		DuelResultWin,
		DuelResultLose,
		DuelResultNA,
	}
}

func GetValidDuelResultStrings() []string {
	return []string{
		DuelResultDraw.String(),
		DuelResultWin.String(),
		DuelResultLose.String(),
		DuelResultNA.String(),
	}
}

func ValidateDuelResultString(s string) DuelResult {
	for _, validDuelResult := range GetValidDuelResults() {
		if s == validDuelResult.String() {
			return validDuelResult
		}
	}
	return DuelResultInvalid
}

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

type PostDuelRequestParam struct {
	LeagueID       NullInt32 `json:"league_id"`
	OpponentUserID string    `json:"opponent_user_id"`
	DeckID         int32     `json:"deck_id"`
	OpponentDeckID int32     `json:"opponent_deck_id"`
	Result         string    `json:"result"`
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

// Create a new request (= unconfirmed duel)
func (c *Controller) PostUserDuelRequest(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("PostUserDuelRequest start")

	// acquire ID of the requested user
	reqUserID := r.Header.Get("UserID")

	// receive body as API parameter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to read body",
			},
			errors.Wrap(err, "")
	}
	var param PostDuelRequestParam
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

	// check if specified IDs exist
	if param.LeagueID.Valid {
		_, err = queries.GetLeague(c.ctx, param.LeagueID.Int32)
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusBadRequest,
				ErrorBody{
					Error: "specified ID not found in the corresponding table: league_id",
				},
				errors.Wrap(err, "")
		} else if err != nil {
			return http.StatusInternalServerError,
				ErrorBody{
					Error: "failed to communicate with database",
				},
				errors.Wrap(err, "")
		}
	}
	_, err = queries.GetUser(c.ctx, param.OpponentUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "specified ID not found in the corresponding table: opponent_user_id",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}
	_, err = queries.GetDeck(c.ctx, param.DeckID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "specified ID not found in the corresponding table: deck_id",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}
	_, err = queries.GetDeck(c.ctx, param.OpponentDeckID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "specified ID not found in the corresponding table: opponent_deck_id",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// check if the given duel result string is valid
	duelResult := ValidateDuelResultString(param.Result)
	if duelResult == DuelResultInvalid {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "specified string invalid: result",
			},
			errors.Wrap(err, "")
	}

	// create request
	err = queries.CreateUnconfirmedDuel(c.ctx, db.CreateUnconfirmedDuelParams{
		LeagueID:  sql.NullInt32(param.LeagueID),
		User1ID:   reqUserID,
		User2ID:   param.OpponentUserID,
		Deck1ID:   param.DeckID,
		Deck2ID:   param.OpponentDeckID,
		Result:    int32(duelResult),
		CreatedBy: 1,
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

	log.Println("PostUserDuelRequest end")

	return http.StatusOK, responseBody, nil
}

// Confirm a duel request
func (c *Controller) PutRequest(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("PutRequest start")

	// acquire URL parameter
	paramDuelIDStr := chi.URLParam(r, "duelId")
	if paramDuelIDStr == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: duelId",
			},
			errors.New("url parameter missing")
	}
	paramDuelID, err := strconv.Atoi(paramDuelIDStr)
	if err != nil {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "failed to parse given url parameter to integer: duelId",
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

	// check if specified ID exists
	_, err = queries.GetDuel(c.ctx, int32(paramDuelID))
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "duel not found (blame: url parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// confirm
	err = queries.ConfirmDuel(c.ctx, int32(paramDuelID))
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

	log.Println("PutRequest end")

	return http.StatusOK, responseBody, nil
}

// Deny / delete a duel request
func (c *Controller) DeleteRequest(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("DeleteRequest start")

	// acquire URL parameter
	paramDuelIDStr := chi.URLParam(r, "duelId")
	if paramDuelIDStr == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: duelId",
			},
			errors.New("url parameter missing")
	}
	paramDuelID, err := strconv.Atoi(paramDuelIDStr)
	if err != nil {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "failed to parse given url parameter to integer: duelId",
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

	// check if specified ID exists
	_, err = queries.GetDuel(c.ctx, int32(paramDuelID))
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "duel not found (blame: url parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// deny / delete
	err = queries.DeleteDuel(c.ctx, int32(paramDuelID))
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

	log.Println("DeleteRequest end")

	return http.StatusOK, responseBody, nil
}
