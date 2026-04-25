package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// User handles requests for user request information.
func User(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		app.Log.Info("serving user info", "user", req)

		response := storage.Owner{
			Address: req.Address,
			Headers: r.Header,
		}

		writeJSON(w, http.StatusOK, response)
	}
}
