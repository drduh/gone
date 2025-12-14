package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/templates"
	"github.com/drduh/gone/version"
)

// Status handles requests for server status and configuration.
func Status(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		app.Log.Info("serving status", "user", req)

		versionInfo := version.Get()
		if !app.ShowBuild {
			versionInfo = map[string]string{
				"id": versionInfo["id"],
			}
		}

		app.CountStorage()
		response := templates.Status{
			Version:  versionInfo,
			Port:     app.Port,
			Uptime:   app.Uptime(),
			Hostname: app.Hostname,
			Index:    app.Index,
			Default:  app.Default,
			Limit:    app.Limit,
			Owner: storage.Owner{
				Address: req.Address,
				Headers: r.Header,
			},
			Sizes: storage.Sizes{
				NumFiles:    app.NumFiles,
				NumMessages: app.NumMessages,
				CharsWall:   app.CharsWall,
			},
		}

		writeJSON(w, http.StatusOK, response)
	}
}
