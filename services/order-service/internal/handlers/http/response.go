package http

import (
	"encoding/json"
	"net/http"
)

// writeJSON пишет JSON-ответ с нужным статусом.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
