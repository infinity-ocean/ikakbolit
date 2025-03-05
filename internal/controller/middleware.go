package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	err string
}

func httpWrapper(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSONtoHTTP(w, http.StatusBadRequest, apiError{err: err.Error()})
		}
	}
}

func writeJSONtoHTTP(w http.ResponseWriter, code int, v any) error {
	log.Println("HTTP to user:", code, v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}