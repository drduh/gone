package handlers

import (
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// Returns content by file name
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if app.Settings.Auth.Require.Download &&
			!auth.Basic(app.Settings.Auth.Header, app.Settings.Auth.Token, r) {
			deny(w, app, req)
			return
		}

		fileName := r.URL.Path[len(app.Settings.Paths.Download):]
		if fileName == "" {
			fileName = r.URL.Query().Get("name")
		}

		var file *config.File
		var found bool

		if fileName != "" {
			for _, rec := range app.Storage.Files {
				if rec.Name == fileName {
					file = rec
					found = true
					break
				}
			}
		}

		if !found {
			writeJSON(w, http.StatusNotFound, responseErrorFileNotFound)
			app.Log.Error(errorFileNotFound,
				"filename", fileName,
				"user", req)
			return
		}

		writeFile(w, file)
		file.Downloads.Total++
		app.Log.Info("served file",
			"filename", file.Name,
			"size", file.Size,
			"downloads", file.Downloads.Total,
			"user", req)

		reason := file.IsExpired(app.Settings)
		if reason != "" {
			app.Storage.Expire(file)
			app.Log.Info("removed file",
				"reason", reason,
				"filename", file.Name,
				"downloads", file.Downloads.Total)
		}
	}
}
