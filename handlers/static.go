package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Static serves pre-defined static (embedded) content from
// templates/data/static.txt in JSON format.
func Static(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		app.Log.Info("serving static content", "user", req)
		writeJSON(w, http.StatusOK, templates.Static{
			Data: templates.StaticData,
		})
	}
}
