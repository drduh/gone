package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// Writes JSON-encoded response
func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

// Returns CSS theme string based on current time of day
func getTheme() string {
	now := time.Now().Hour()
	if now > 7 && now < 19 {
		return "light"
	}
	return "dark"
}
