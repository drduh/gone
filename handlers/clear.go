package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Clear handles requests to clear Storage content.
func Clear(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if !app.Allow(app.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests,
				errorJSON(app.RateLimit))
			app.Log.Error(app.RateLimit, "user", req)
			return
		}

		app.ClearStorage()
		app.CountStorage()
		app.Log.Info("storage cleared", "user", req)

		writeJSON(w, http.StatusOK, app.Sizes)
	}
}
