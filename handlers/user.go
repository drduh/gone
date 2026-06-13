package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/templates"
)

// UserInfo handles requests for user request information.
func UserInfo(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		app.Log.Info("serving user request info",
			"user", req)

		response := templates.User{
			Owner: storage.Owner{
				Address: req.Address,
				Mask:    req.Mask,
				Headers: r.Header,
			},
			IsBrowser: req.IsBrowser,
		}

		writeJSON(w, http.StatusOK, response)
	}
}
