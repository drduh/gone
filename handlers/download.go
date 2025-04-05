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

		fileName := getFilename(r, len(app.Download))
		if fileName == "" {
			writeJSON(w, http.StatusNotFound, errorJSON(app.NoFilename))
			app.Log.Error(app.NoFilename, "user", req)
			return
		}

		file := app.FindFile(fileName)
		if file == nil {
			writeJSON(w, http.StatusNotFound, errorJSON(app.NotFound))
			app.Log.Error(app.NotFound, "filename", fileName, "user", req)
			return
		}

		file.Serve(w)
		app.Log.Info("served file",
			"filename", file.Name, "size", file.Size,
			"downloads", file.Total, "user", req)

		reason := file.IsExpired(app.Settings)
		if reason != "" {
			app.Expire(file)
			app.Log.Info("removed file",
				"reason", reason, "filename", file.Name,
				"downloads", file.Total)
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
