package handlers

import (
	"encoding/json"
	"net/http"
)

// Server Heartbeat JSON response
func Heartbeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()
		payload := map[string]interface{}{
			"ip": ip,
			"ua": ua,
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(payload)
	}
}
