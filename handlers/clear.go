package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Clear handles requests to clear Storage content.
func Clear(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		app.ClearStorage()
		app.Log.Info("storage cleared",
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
		} else {
			writeJSON(w, http.StatusOK, "storage cleared")
		}
	}
}
