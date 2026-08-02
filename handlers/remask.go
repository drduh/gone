package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// UserRemask assigns a new address mask to the user.
func UserRemask(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		refreshMask(req.Address)
		reqNew := parseRequest(r)
		app.Log.Info("remasked user",
			"new", reqNew.AddressMask,
			"old", req.AddressMask)

		if req.IsBrowser {
			toPath(w, r, app.Root)
		} else {
			writeJSON(w, http.StatusOK,
				req.AddressMask+" is now "+reqNew.AddressMask)
		}
	}
}
