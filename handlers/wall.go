package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Wall handles wall form submissions and updates content in Storage.
func Wall(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if r.Method == http.MethodPost {
			if r.FormValue("clear") != "" {
				app.Log.Debug("clearing wall",
					"length", app.CountWall(), "user", req)
				app.ClearWall()
			}

			wallContent := r.FormValue("wall")
			if wallContent != "" {
				app.WallContent = wallContent
				app.Log.Debug("updated wall content",
					"length", app.CountWall(), "user", req)
			}

			toRoot(w, r, app.Root)
		}

		writeJSON(w, http.StatusOK, app.WallContent)
	}
}
