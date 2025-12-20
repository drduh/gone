package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Static handles requests for static embedded content.
func Static(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		app.Log.Info("serving static", "user", req)
		writeJSON(w, http.StatusOK, templates.Static{
			Data: templates.StaticData,
		})
	}
}
