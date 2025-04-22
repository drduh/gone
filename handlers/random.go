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
		app.Log.Info("serving random", "user", req)

		var response string
		path := getParam(r, len(app.Random), "random")
		switch path {
		case "name":
			response = util.RandomName()
		case "number":
			response = util.RandomNumber()
		default:
			response = util.RandomName() + util.RandomNumber()
		}

		writeJSON(w, http.StatusOK, response)
	}
}
