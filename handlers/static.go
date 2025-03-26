package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Static content JSON response
func Static(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		writeJSON(w, http.StatusOK, templates.Static{
			Data: templates.StaticData,
		})
		app.Log.Info("served static content", "user", req)
	}
}
