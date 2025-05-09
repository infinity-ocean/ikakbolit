package reply

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, v any) error {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
