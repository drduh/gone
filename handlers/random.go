package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// Random handles requests to generate a random string.
func Random(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		path := getRequestParameter(r, len(app.Random), "random")
		app.Log.Info("serving random", "user", req, "path", path)
		writeJSON(w, http.StatusOK, util.GetRandom(path))
	}
}
