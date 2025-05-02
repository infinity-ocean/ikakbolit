package controller

import (
	"net/http"
	"net/url"
	"strings"
)

var SensitiveQueryParams = map[string]struct{}{
	"user_id": {},
}

var SensitiveHeaders = map[string]struct{}{
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
		if _, bad := SensitiveQueryParams[strings.ToLower(key)]; bad {
			continue
		}
		safe[key] = val
	}
	return safe.Encode()
}

func sanitizeHeaders(headers http.Header) map[string]string {
	safe := make(map[string]string, len(headers))
	for k, vals := range headers {
		if _, bad := SensitiveHeaders[strings.ToLower(k)]; bad {
			continue
		}
		safe[k] = strings.Join(vals, ", ")
	}
	return safe
}