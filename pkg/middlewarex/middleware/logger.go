package middleware

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/go-chi/chi/v5/middleware"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func LoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			headers := sanitizeHeaders(r.Header)

			req := log.With(
				slog.Time("received_at", t1),
				slog.String("method", r.Method),
				slog.String("url", r.URL.Path),
				slog.String("query", sanitizeQuery(r.URL.Query())),
				slog.Any("headers", headers),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("client_ip", r.RemoteAddr),
				slog.String(string(requestIDKey), middleware.GetReqID(r.Context())),
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

var sensitiveQueryParams = map[string]struct{}{
	"user_id": {},
}

var sensitiveHeaders = map[string]struct{}{
	"authorization":       {},
	"proxy-authorization": {},
	"cookie":              {},
	"set-cookie":          {},
	"x-api-key":           {},
	"x-csrf-token":        {},
	"x-xsrf-token":        {},
	"x-forwarded-for":     {},
	"forwarded":           {},
	"referer":             {},
	"referrer-policy":     {},
	"from":                {},
}

func sanitizeQuery(values url.Values) string {
	safe := url.Values{}
	for key, val := range values {
		if _, bad := sensitiveQueryParams[strings.ToLower(key)]; bad {
			continue
		}
		safe[key] = val
	}
	return safe.Encode()
}

func sanitizeHeaders(headers http.Header) map[string]string {
	safe := make(map[string]string, len(headers))
	for k, vals := range headers {
		if _, bad := sensitiveHeaders[strings.ToLower(k)]; bad {
			continue
		}
		safe[k] = strings.Join(vals, ", ")
	}
	return safe
}
