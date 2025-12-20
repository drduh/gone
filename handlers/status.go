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
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
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
			Sizes:    app.Sizes,
			Owner: storage.Owner{
				Address: req.Address,
				Headers: r.Header,
			},
		}

		writeJSON(w, http.StatusOK, response)
	}
}
