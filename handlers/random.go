package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// Random serves a random string of specified type.
func Random(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		path := getParam(r, len(app.Random), "random")
		app.Log.Info("serving random", "user", req, "path", path)
		writeJSON(w, http.StatusOK, util.GetRandom(path))
	}
}
