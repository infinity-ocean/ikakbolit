package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

func httpWrapper(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch {
		case errors.Is(err, model.ErrNoContent):
			writeJSONtoHTTP(w, http.StatusNoContent, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrBadRequest):
			writeJSONtoHTTP(w, http.StatusBadRequest, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrNotFound):
			writeJSONtoHTTP(w, http.StatusNotFound, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrMethodNotAllowed):
			writeJSONtoHTTP(w, http.StatusMethodNotAllowed, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrInternalServerError):
			writeJSONtoHTTP(w, http.StatusInternalServerError, APIError{Message: err.Error()})
		default:
			writeJSONtoHTTP(w, http.StatusInternalServerError, APIError{Message: err.Error()})
		}
	})
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
