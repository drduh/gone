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

		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusOK, app.WallContent)
			return
		}

		if r.PostFormValue(formFieldClear) != "" {
			app.Log.Debug("clearing wall",
				"length", app.WallChars,
				"user", req)
			app.ClearWall()
			app.Log.Info("cleared wall",
				"user", req)
		}

		formContent := r.PostFormValue(formFieldWall)
		if formContent != "" {
			length := charCount(formContent)
			if length > app.WallLimits.LengthChars {
				writeJSON(w, http.StatusBadRequest,
					errorJSON(app.WallLimit))
				app.Log.Error(app.WallLimit,
					"limit", app.WallLimits.LengthChars,
					"length", length,
					"user", req)
				return
			}

			now := time.Now()
			app.WallModifiedTime = now
			app.WallModifiedTimeFmt = now.Format(app.TimeFormat)
			app.WallModifiedUser = req.AddressMask

			app.Log.Debug("updating wall",
				"length", length,
				"user", req)
			app.WallContent = formContent
			app.CountWall()
			app.Log.Info("updated wall",
				"length", length,
				"user", req)
		}

		if r.PostFormValue(formFieldDownload) == formFieldWall {
			app.ServeWall(w)
			app.Log.Info("downloaded wall",
				"user", req)
			return
		}

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, app.WallContent)
	}
}
