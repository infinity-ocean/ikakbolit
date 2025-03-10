package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type APIError struct {
	Err string `json:"error"`
}

func httpWrapper(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			errString := err.Error()

			switch {
			case errString == "no upcoming schedules found: all schedules are expired or user has no schedules":
				w.WriteHeader(http.StatusNoContent)
				return
			case strings.HasPrefix(errString, "failed to insert schedule"):
				writeJSONtoHTTP(w, http.StatusInternalServerError, APIError{Err: errString})
			default:
				writeJSONtoHTTP(w, http.StatusBadRequest, APIError{Err: errString})
			}
		}
	}
}

func writeJSONtoHTTP(w http.ResponseWriter, code int, v any) error {
	log.Println("HTTP to user:", code, v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		return json.NewEncoder(w).Encode(v)
	}
	return nil
}
