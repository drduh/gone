package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Wall handles requests to read and modify Wall content in Storage.
func Wall(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		app.CountWall()

		if r.Method == http.MethodPost {
			if r.FormValue(formFieldClear) != "" {
				app.Log.Debug("clearing wall",
					"length", app.CharsWall,
					"user", req)
				app.ClearWall()
				app.Log.Info("cleared wall",
					"user", req)
			}

			formContent := r.FormValue(formFieldWall)
			if formContent != "" {
				app.Log.Debug("updating wall",
					"length", len(formContent),
					"user", req)

				now := time.Now()
				app.WallModified = now
				app.WallModifiedFmt = now.Format(app.TimeFormat)
				app.WallContent = formContent

				app.Log.Info("updated wall",
					"length", len(app.WallContent),
					"user", req)
			}

			formContent = r.FormValue("download")
			if formContent == "wall" {
				app.ServeWall(w)
				app.Log.Info("downloaded wall",
					"user", req)
				return
			}

			if req.IsBrowser {
				toPath(w, r, app.Root)
				return
			}
		}

		writeJSON(w, http.StatusOK, app.WallContent)
	}
}
