package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
	"github.com/drduh/gone/version"
)

// Returns server status response
func Heartbeat(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		uptime := app.Uptime()
		app.Log.Info("serving heartbeat",
			"uptime", uptime, "user", req)

		response := templates.Heartbeat{
			Hostname:     app.Hostname,
			Version:      version.Short(),
			Port:         app.Settings.Port,
			Uptime:       uptime,
			FileCount:    app.Storage.CountFiles(),
			MessageCount: app.Storage.CountMessages(),
			Limits:       app.Settings.Limits,
			Owner: config.Owner{
				Address: req.Address,
				Headers: r.Header,
			},
		}

		writeJSON(w, http.StatusOK, response)
	}
}
