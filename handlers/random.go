package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// Random handles requests for random strings.
func Random(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		count := app.RandomLimits.StrCount
		path := getRequestParameter(
			r, len(app.Random), "random")

		response := make([]string, count)
		for i := range count {
			response[i] = util.GetRandom(path)
		}

		app.Log.Info("serving random", "user", req)

		if req.IsBrowser {
			renderIndex(app, w, r, req, path, response)
		} else {
			writeJSON(w, http.StatusOK, response)
		}
	}
}
