package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tsubasa283paris/HMDC/sqlc/db"
	"github.com/tsubasa283paris/HMDC/utils"

	"github.com/pkg/errors"
)

type HelloParam struct {
	Param1 string `json:"param1"`
	Param2 int64  `json:"param2"`
}

type HelloResponse struct {
	Message  string `json:"message"`
	NumUsers int    `json:"numUsers"`
}

// Concat received parameters
func (c *Controller) Hello(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	log.Println("Hello start")

	// receive body as API parameter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to read body",
			},
			errors.Wrap(err, "")
	}
	var param HelloParam
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
	userList, err := queries.ListUsers(c.ctx)
	if err != nil {
		return http.StatusInternalServerError,
			ErrorBody{
				Error: "failed to communicate with database",
			},
			errors.Wrap(err, "")
	}

	// write response
	resp := HelloResponse{
		Message:  "Hello " + param.Param1 + strconv.FormatInt(param.Param2, 10) + "!",
		NumUsers: len(userList),
	}

	log.Println("Hello end")

	return http.StatusOK,
		resp,
		nil
}
