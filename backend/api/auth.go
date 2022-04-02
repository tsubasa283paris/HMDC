package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type LoginParam struct {
	id       string `json:"id"`
	password string `json:"password"`
}

type LoginResponse struct {
	token string `json:"token"`
}

// Check user id and password, return a token if valid
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login start")

	// receive body as API parameter
	var param LoginParam
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		RespondError(
			w,
			"failed to decode body string to JSON format required by this API",
			http.StatusBadRequest,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}
	log.Println("param:", param)

	// open database connection
	dbCnx, err := utils.DbCnx()
	if err != nil {
		RespondError(
			w,
			"failed to connect to the database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// prepare for query
	queries := db.New(dbCnx)

	// run query
	user, err := queries.GetUser(h.ctx, param.id)
	if errors.Is(err, sql.ErrNoRows) {
		RespondError(
			w,
			"invalid id or password",
			http.StatusBadRequest,
		)
		log.Println("Invalid ID or password")
		return
	} else if err != nil {
		RespondError(
			w,
			"failed to communicate with database",
			http.StatusInternalServerError,
		)
		log.Println(fmt.Sprintf("%+v", errors.Wrap(err, "")))
		return
	}

	// invalid if password is wrong
	if param.password != user.Password {
		RespondError(
			w,
			"invalid id or password",
			http.StatusBadRequest,
		)
		log.Println("Invalid ID or password")
		return
	}

	// write response
	resp := LoginResponse{
		token: "admin",
	}
	respondJSON(w, resp)

	log.Println("Login end")
}
