package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Server Heartbeat JSON response
func Heartbeat(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()
		payload := map[string]interface{}{
			"ip":     ip,
			"ua":     ua,
			"uptime": time.Since(app.Start).String(),
		}
		writeJSON(w, http.StatusOK, payload)
	}
}
