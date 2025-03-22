package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// Writes JSON-encoded response
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// Writes data response for Files
func writeData(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// Returns CSS theme based on current time of day if unset
func getTheme(theme string) string {
	if theme != "" {
		return theme
	}
	now := time.Now().Hour()
	if now > 7 && now < 19 {
		return "light"
	}
	return "dark"
}
