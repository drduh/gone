package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Clear handles requests to clear Storage content.
func Clear(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		app.ClearStorage()
		app.Log.Info("storage cleared", "user", req)
		if req.IsBrowser {
			toRoot(w, r, app.Root)
		}
	}
}
