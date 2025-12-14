package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Wall handles requests to read and modify Wall content in Storage.
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
				app.Log.Info("cleared wall", "user", req)
			}

			wallContent := r.FormValue("wall")
			if wallContent != "" {
				app.Log.Debug("updating wall",
					"length", app.CountWall(), "user", req)
				app.WallContent = wallContent
				app.Log.Info("updated wall", "user", req)
			}

			toRoot(w, r, app.Root)
		}

		writeJSON(w, http.StatusOK, app.WallContent)
	}
}
