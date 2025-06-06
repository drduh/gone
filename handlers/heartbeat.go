package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/templates"
	"github.com/drduh/gone/version"
)

// Heartbeat handles requests for server status/configuration.
func Heartbeat(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		uptime := app.Uptime()
		app.Log.Info("serving heartbeat",
			"uptime", uptime, "user", req)

		response := templates.Heartbeat{
			Hostname:     app.Hostname,
			Index:        app.Index,
			Limits:       app.Limits,
			FileCount:    app.CountFiles(),
			MessageCount: app.CountMessages(),
			Owner: storage.Owner{
				Address: req.Address,
				Headers: r.Header,
			},
			Port:    app.Port,
			Uptime:  uptime,
			Version: version.Get(),
		}

		writeJSON(w, http.StatusOK, response)
	}
}
