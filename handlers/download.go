package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Download handles requests to download a File from Storage by name.
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
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
		app.Log.Debug("file found", "filename", file.Name, "user", req)

		file.Serve(w)
		app.Log.Info("served file",
			"filename", file.Name, "size", file.Size,
			"downloads", file.Total, "user", req)

		reason := file.IsExpired()
		if reason != "" {
			app.Expire(file)
			app.Log.Info("removed file",
				"reason", reason, "filename", file.Name,
				"downloads", file.Total)
		}
	}
}
