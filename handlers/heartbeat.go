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

		response := templates.Heartbeat{
			Hostname:  app.Hostname,
			Version:   app.Version,
			Port:      app.Settings.Port,
			Uptime:    uptime,
			FileCount: len(app.Storage.Files),
			Limits: templates.Limits{
				Downloads: app.Settings.Limits.Downloads,
			},
			Client: templates.Client{
				Address: ip,
				Headers: r.Header,
			},
		}

		writeJSON(w, http.StatusOK, response)
		app.Log.Info("heartbeat",
			"uptime", uptime,
			"ip", ip, "ua", ua)

	}
}
