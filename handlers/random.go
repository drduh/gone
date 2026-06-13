package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// Random handles requests for a random string.
func Random(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		path := getRequestParameter(r,
			len(app.Random), "random")
		count := app.RandomLimits.StrCount

		results := make([]string, count)
		for i := range count {
			results[i] = util.GetRandom(path)
		}

		app.Log.Info("serving random",
			"count", count,
			"path", path,
			"user", req)

		writeJSON(w, http.StatusOK, results)
	}
}
