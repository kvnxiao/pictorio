package response

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, value interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}
