package api

import (
	"database/sql"
	"encoding/json"
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
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("Login start")

	// receive body as API parameter
	var param LoginParam
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to decode body string to JSON format required by this API",
			},
			errors.Wrap(err, "")
	}
	log.Println("param:", param)

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
	user, err := queries.GetUser(c.ctx, param.id)
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "invalid id or password",
			},
			errors.Wrap(err, "")
	} else if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// invalid if password is wrong
	if param.password != user.Password {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "invalid id or password",
			},
			errors.Wrap(err, "")
	}

	// write response
	resp := LoginResponse{
		token: "admin",
	}

	log.Println("Login end")

	return http.StatusOK,
		resp,
		nil
}
