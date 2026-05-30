package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// UserRemask assigns a new address mask to the user.
func UserRemask(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		refreshMask(req.Address)
		reqNew := parseRequest(r)
		app.Log.Info("remasked user",
			"new", reqNew.Mask,
			"old", req.Mask)

		if req.IsBrowser {
			toPath(w, r, app.Root)
		} else {
			writeJSON(w, http.StatusOK,
				req.Mask+" is now "+reqNew.Mask)
		}
	}
}
