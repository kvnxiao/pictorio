package response

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
