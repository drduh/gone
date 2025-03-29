package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Returns content by file name
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if !isAllowed(app, r) {
			deny(w, app, req)
			return
		}

		fileName := getFilename(r, len(app.Settings.Paths.Download))
		if fileName == "" {
			writeJSON(w, http.StatusNotFound, errorJSON(app.Error.NoFilename))
			app.Log.Error(app.Error.NoFilename, "user", req)
			return
		}

		var file *config.File
		var found bool

		for _, rec := range app.Storage.Files {
			if rec.Name == fileName {
				file = rec
				found = true
				break
			}
		}

		if !found {
			writeJSON(w, http.StatusNotFound, errorJSON(app.Error.NotFound))
			app.Log.Error(app.Error.NotFound,
				"filename", fileName, "user", req)
			return
		}

		writeFile(w, file)
		file.Downloads.Total++
		app.Log.Info("served file",
			"filename", file.Name, "size", file.Size,
			"downloads", file.Downloads.Total, "user", req)

		reason := file.IsExpired(app.Settings)
		if reason != "" {
			app.Storage.Expire(file)
			app.Log.Info("removed file",
				"reason", reason, "filename", file.Name,
				"downloads", file.Downloads.Total)
		}
	}
}

// Returns filename value from request
func getFilename(r *http.Request, pathLen int) string {
	f := r.URL.Path[pathLen:]
	if f == "" {
		f = r.URL.Query().Get("name")
	}
	if f == "" {
		f = r.FormValue("name")
	}
	return f
}
