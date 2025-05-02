package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/infinity-ocean/ikakbolit/internal/model"
)

func LoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			traceID := r.Header.Get("X-TRACE-ID")
			if traceID == "" {
				traceID = strconv.FormatInt(time.Now().UnixNano(), 36)
			}

			headers := sanitizeHeaders(r.Header)

			req := log.With(
				slog.Time("received_at", t1),
				slog.String("trace_id", traceID),
				slog.String("method", r.Method),
				slog.String("url", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.Any("headers", headers),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("client_ip", r.RemoteAddr),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
					req.Info("request completed",
					slog.String("send_at", time.Now().String()),
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func httpWrapper(f handlerWithError) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch {
		case errors.Is(err, model.ErrNoContent):
			sendJSONtoHTTP(w, http.StatusNoContent, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrBadRequest):
			sendJSONtoHTTP(w, http.StatusBadRequest, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrNotFound):
			sendJSONtoHTTP(w, http.StatusNotFound, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrMethodNotAllowed):
			sendJSONtoHTTP(w, http.StatusMethodNotAllowed, APIError{Message: err.Error()})
		case errors.Is(err, model.ErrInternalServerError):
			sendJSONtoHTTP(w, http.StatusInternalServerError, APIError{Message: err.Error()})
		default:
			sendJSONtoHTTP(w, http.StatusInternalServerError, APIError{Message: err.Error()})
		}
	})
}	

func sendJSONtoHTTP(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		return json.NewEncoder(w).Encode(v)
	}
	return nil
}

func sanitizeHeaders(headers http.Header) map[string]string {
	safe := make(map[string]string)
	for k, v := range headers {
		if strings.EqualFold(k, "user_id") {
			continue
		}
		safe[k] = strings.Join(v, ", ")
	}
	return safe
}