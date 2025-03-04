package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Static content JSON response
func Static(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()
		resp := templates.Static{
			Data: templates.StaticData,
		}
		writeJSON(w, http.StatusOK, resp)
		app.Log.Info("served static content",
			"ip", ip, "ua", ua)
	}
}
