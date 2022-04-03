package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type UserStatsPerLeague struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

type UserDuelHistory struct {
	DuelID         int32     `json:"duel_id"`
	LeagueID       NullInt32 `json:"league_id"`
	OpponentUserID string    `json:"opponent_user_id"`
	DeckID         int32     `json:"deck_id"`
	OpponentDeckID int32     `json:"opponent_deck_id"`
	Result         string    `json:"result"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserDeck struct {
	DeckID      int32     `json:"deck_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LeagueID    NullInt32 `json:"league_id"`
	NumDuel     int32     `json:"num_duel"`
	NumWin      int32     `json:"num_win"`
}

type UserDetails struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type PutUserDetailsParam struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
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
				Error: "url parameter missing: userId",
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
				Error: "user not found (blame: url parameter)",
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
	statsList, err := queries.ListUserStats(c.ctx, paramUserID)
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

// Get duel history of the specified user
func (c *Controller) GetUserDuelHistory(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUserDuelHistory start")

	// acquire URL parameter
	paramUserID := chi.URLParam(r, "userId")
	if paramUserID == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: userId",
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
				Error: "user not found (blame: url parameter)",
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
	duelHistoryList, err := queries.ListUserDuelHistory(c.ctx, paramUserID)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []UserDuelHistory
	for _, duelHistory := range duelHistoryList {
		responseBody = append(responseBody, UserDuelHistory{
			DuelID:         duelHistory.ID,
			LeagueID:       NullInt32(duelHistory.LeagueID),
			OpponentUserID: duelHistory.OpponentUserID,
			DeckID:         duelHistory.DeckID,
			OpponentDeckID: duelHistory.OpponentDeckID,
			Result:         fmt.Sprintf("%v", duelHistory.Result),
			CreatedAt:      duelHistory.CreatedAt,
		})
	}

	log.Println("GetUserDuelHistory end")

	return http.StatusOK, responseBody, nil
}

// Get decks belonging to the specified user
func (c *Controller) GetUserDecks(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUserDecks start")

	// acquire URL parameter
	paramUserID := chi.URLParam(r, "userId")
	if paramUserID == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: userId",
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
				Error: "user not found (blame: url parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// get decks
	deckList, err := queries.ListUserDecks(c.ctx, paramUserID)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// create responseBody
	var responseBody []UserDeck
	for _, deck := range deckList {
		responseBody = append(responseBody, UserDeck{
			DeckID:      deck.ID,
			Name:        deck.Name,
			Description: deck.Description,
			LeagueID:    NullInt32(deck.LeagueID),
			NumDuel:     deck.NumDuel,
			NumWin:      deck.NumWin,
		})
	}

	log.Println("GetUserDecks end")

	return http.StatusOK, responseBody, nil
}

// Get user details containing name and password
// Hide password if not requested by the target user itself
func (c *Controller) GetUserDetails(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("GetUserDetails start")

	// acquire URL parameter
	paramUserID := chi.URLParam(r, "userId")
	if paramUserID == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: userId",
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

	// get user details
	user, err := queries.GetUser(c.ctx, paramUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "user not found (blame: url parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// hide password if not self-querying
	password := ""
	reqUserID := r.Header.Get("UserID")
	if reqUserID != paramUserID {
		password = "*"
	} else {
		password = user.Password
	}

	// create responseBody
	responseBody := UserDetails{
		Name:     user.Name,
		Password: password,
	}

	log.Println("GetUserDetails end")

	return http.StatusOK, responseBody, nil
}

// Edit user details: id, name and password
// Return 403 if not requested by the target user itself
func (c *Controller) PutUserDetails(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("PutUserDetails start")

	// acquire URL parameter
	paramUserID := chi.URLParam(r, "userId")
	if paramUserID == "" {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "url parameter missing: userId",
			},
			nil
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
	var param PutUserDetailsParam
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

	// check if the specified ID exists
	_, err = queries.GetUser(c.ctx, paramUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound,
			ErrorBody{
				Error: "user not found (blame: url parameter)",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// reject if not self-querying
	reqUserID := r.Header.Get("UserID")
	if reqUserID != paramUserID {
		return http.StatusForbidden,
			ErrorBody{
				Error: "updating other user is not allowed",
			},
			errors.New("forbidden user update")
	}

	// update user details
	err = queries.UpdateUser(c.ctx, db.UpdateUserParams{
		ID:       paramUserID,
		ID_2:     param.ID,
		Password: param.Password,
		Name:     param.Name,
	})
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(err.Code)
		if err.Code == "unique_violation" {
			return http.StatusBadRequest,
				ErrorBody{
					Error: "specified id already exists",
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

	// create responseBody
	responseBody := ErrorBody{
		Error: "",
	}

	log.Println("PutUserDetails end")

	return http.StatusOK, responseBody, nil
}
