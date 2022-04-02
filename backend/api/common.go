package api

import (
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
func errorResponse(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(errorBodyJSONBytes(err))
}

// Marshal NullInt32 value converting null value as JSON null
func (x NullInt32) MarshalJSON() ([]byte, error) {
	if x.Valid {
		return json.Marshal(x.Int32)
	}
	return []byte(`null`), nil
}
