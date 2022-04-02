package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error string `json:"error"`
}

type NullInt32 struct {
	sql.NullInt32
}

// Handler struct
type Handler struct {
	ctx context.Context
}

// Constructor for struct Handler
func NewHandler() *Handler {
	return &Handler{
		ctx: context.Background(),
	}
}

// Generate and write JSON for response as success
func respondJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		RespondError(
			w,
			"failed to json.Marshal",
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// Convert error string to JSON form byte array
func errorBodyJSONBytes(s string) []byte {
	errorBodyJSON, err := json.Marshal(ErrorBody{
		Error: s,
	})
	if err != nil {
		errorBodyJSON = []byte{}
	}
	return errorBodyJSON
}

// Write error code to header and error JSON to body respectively
func RespondError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(errorBodyJSONBytes(err))
}

// Marshal NullInt32 value converting null value as JSON null
func (x NullInt32) MarshalJSON() ([]byte, error) {
	if x.Valid {
		return json.Marshal(x.Int32)
	}
	return []byte(`null`), nil
}
