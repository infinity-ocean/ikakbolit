package reply

import (
	"errors"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/pkg/errcodes"
	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func ErrorDecorator(f HandlerWithError) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch {
		case errors.Is(err, errcodes.ErrNotFound):
			JSON(w, http.StatusNotFound, rest.APIError{
				ErrorCode: "404",
				Message: "resource not found",
			})
		case errors.Is(err, errcodes.ErrBadRequest):
			JSON(w, http.StatusBadRequest, rest.APIError{
				ErrorCode: "400",
				Message: "invalid input data",
			})
		default:
			JSON(w, http.StatusInternalServerError, rest.APIError{
				ErrorCode: "500",
				Message: "internal server error",
			})
		}
	})
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error
