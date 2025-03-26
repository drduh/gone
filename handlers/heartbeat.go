package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Returns server status response
func Heartbeat(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		uptime := time.Since(app.Start).String()

		response := templates.Heartbeat{
			Hostname:  app.Hostname,
			Version:   app.Version,
			Port:      app.Settings.Port,
			Uptime:    uptime,
			FileCount: app.Storage.CountFiles(),
			Limits:    app.Settings.Limits,
			Owner: config.Owner{
				Address: req.Address,
				Headers: r.Header,
			},
		}

		writeJSON(w, http.StatusOK, response)
		app.Log.Info("served heartbeat",
			"uptime", uptime, "user", req)
	}
}
