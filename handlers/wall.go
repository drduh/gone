package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Wall handles requests to read and modify Wall content.
func Wall(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		if r.Method != http.MethodPost {
			content, _ := app.GetWall()
			writeJSON(w, http.StatusOK, content)
			return
		}

		if r.PostFormValue(formFieldDownload) == formFieldWall {
			_, meta := app.GetWall()
			app.ServeWall(w)
			app.Log.Info("downloaded wall",
				"length", meta.WallChars,
				"user", req)
			return
		}

		if r.PostFormValue(formFieldClear) != "" {
			app.ClearWall()
			app.Log.Info("cleared wall",
				"user", req)
		}

		formContent := r.PostFormValue(formFieldWall)
		length := charCount(formContent)

		if length > app.WallLimits.LengthChars {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.WallLimit))
			app.Log.Error(app.WallLimit,
				"excess", length-app.WallLimits.LengthChars,
				"length", length,
				"limit", app.WallLimits.LengthChars,
				"user", req)
			return
		}

		now := time.Now()
		meta := app.SetWall(formContent, storage.WallMeta{
			WallModifiedTime:    now,
			WallModifiedTimeFmt: now.Format(app.TimeFormat),
			WallModifiedUser:    req.AddressMask,
		}, app.WallLimits.LengthChars)

		app.Log.Info("updated wall",
			"capacity", meta.WallCap,
			"length", length,
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, formContent)
	}
}
