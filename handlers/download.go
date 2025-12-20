package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Download handles requests to download a File from Storage.
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		fileName := getParam(r, len(app.Download), "name")
		if fileName == "" {
			writeJSON(w, http.StatusNotFound, errorJSON(app.NoFilename))
			app.Log.Error(app.NoFilename, "user", req)
			return
		}
		app.Log.Debug("file requested", "filename", fileName, "user", req)

		file := app.FindFile(fileName)
		if file == nil {
			writeJSON(w, http.StatusNotFound, errorJSON(app.NotFound))
			app.Log.Error(app.NotFound, "filename", fileName, "user", req)
			return
		}

		file.Serve(w)
		app.Log.Info("served file", "id", file.Id, "name", file.Name,
			"size", file.Size, "downloads", file.Total, "user", req)

		reason := file.IsExpired()
		if reason != "" {
			app.Expire(file)
			app.Log.Info("removed file", "reason", reason,
				"id", file.Id, "name", file.Name, "downloads", file.Total)
		}
	}
}
