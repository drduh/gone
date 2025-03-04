package handlers

import (
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
		writeJSON(w, http.StatusOK, payload)
	}
}
