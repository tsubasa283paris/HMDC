package api

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type LoginParam struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SignUpParam struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Check user id and password, return a token if valid
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("Login start")

	// receive body as API parameter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to read body",
			},
			errors.Wrap(err, "")
	}
	var param LoginParam
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

	// run query
	user, err := queries.GetUser(c.ctx, param.ID)
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
	if param.Password != user.Password {
		return http.StatusBadRequest,
			ErrorBody{
				Error: "invalid id or password",
			},
			errors.Wrap(err, "")
	}

	// write response
	resp := LoginResponse{
		Token: "admin",
	}

	log.Println("Login end")

	return http.StatusOK,
		resp,
		nil
}

// Check user id and password, return OK and register to the database
// if user id doesn't match any existings
func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("SignUp start")

	// receive body as API parameter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to read body",
			},
			errors.Wrap(err, "")
	}
	var param SignUpParam
	err = json.Unmarshal(body, &param)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to decode body string to JSON format required by this API",
			},
			errors.Wrap(err, "")
	}
	log.Println("param:", string(body))

	// open database connection and transaction
	dbCnx, err := utils.DbCnx()
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to connect to the database",
			},
			errors.Wrap(err, "")
	}
	dbTx, err := dbCnx.Begin()
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to begin transaction with the database",
			},
			errors.Wrap(err, "")
	}

	// prepare for query
	queries := db.New(dbCnx)

	// run query
	_, err = queries.GetUser(c.ctx, param.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return http.StatusInternalServerError,
				ErrorBody{
					Error: "failed to communicate with database",
				},
				errors.Wrap(err, "")
		} else {
			return http.StatusBadRequest,
				ErrorBody{
					Error: "specified id already exists",
				},
				errors.Wrap(err, "")
		}
	}

	// register to the database
	_, err = queries.CreateUser(c.ctx, db.CreateUserParams{
		ID:       param.ID,
		Password: param.Password,
		Name:     param.Name,
	})
	if err != nil {
		_ = dbTx.Rollback()
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// commit transaction
	err = dbTx.Commit()
	if err != nil {
		_ = dbTx.Rollback()
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to commit transaction to the database",
			},
			errors.Wrap(err, "")
	}

	// write response
	resp := ErrorBody{
		Error: "",
	}

	log.Println("SignUp end")

	return http.StatusOK,
		resp,
		nil
}
