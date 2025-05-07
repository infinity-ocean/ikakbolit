package reply

import (
	"errors"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/internal/dto"
	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func ErrorDecorator(f HandlerWithError) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch {
		case errors.Is(err, dto.ErrNotFound):
			JSON(w, http.StatusNotFound, rest.APIError{Message: "resource not found"})
		case errors.Is(err, dto.ErrBadRequest):
			JSON(w, http.StatusBadRequest, rest.APIError{Message: "invalid input data"})
		default:
			JSON(w, http.StatusInternalServerError, rest.APIError{Message: "internal server error"})
		}
	})
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error