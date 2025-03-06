package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Server Heartbeat JSON response
func Heartbeat(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, ua := r.RemoteAddr, r.UserAgent()
		uptime := time.Since(app.Start).String()

		resp := templates.Heartbeat{
			Hostname:  app.Hostname,
			Version:   app.Version,
			Port:      app.Settings.Port,
			Uptime:    uptime,
			FileCount: len(app.Storage.Files),
		}

		writeJSON(w, http.StatusOK, resp)
		app.Log.Info("heartbeat",
			"uptime", uptime,
			"ip", ip, "ua", ua)

	}
}
