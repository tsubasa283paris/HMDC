package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type ErrorBody struct {
	Error string `json:"error"`
}

type NullInt32 sql.NullInt32

// Controller struct
type Controller struct {
	ctx context.Context
}

type Handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

// Constructor for struct Controller
func NewController() *Controller {
	return &Controller{
		ctx: context.Background(),
	}
}

// Generate and write JSON for response as success
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to marshal response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rv := recover(); rv != nil {
			debug.PrintStack()
			log.Printf("panic: %s", rv)
			http.Error(w, http.StatusText(
				http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	status, res, err := h(w, r)
	if err != nil {
		log.Printf("error: %s", err)
	}
	RespondJSON(w, status, res)
	return
}

// Marshal NullInt32 value converting null value as JSON null
func (x NullInt32) MarshalJSON() ([]byte, error) {
	if x.Valid {
		return json.Marshal(x.Int32)
	}
	return []byte(`null`), nil
}
