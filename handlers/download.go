package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Download handles requests to download a File from Storage.
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		filename := getRequestParameter(r, len(app.Download), "name")
		if filename == "" {
			app.Log.Error(app.NoFilename, "user", req)
			writeJSON(w, http.StatusNotFound, errorJSON(app.NoFilename))
			return
		}

		file, reason := app.ReserveDownload(filename)
		if file == nil {
			app.Log.Error(app.NotFound,
				"reason", reason,
				"user", req)
			writeJSON(w, http.StatusNotFound, errorJSON(app.NotFound))
			return
		}
		app.Log.Debug("file found",
			"file", file.Name,
			"user", req)

		if app.Write.Duration > 0 {
			_ = http.NewResponseController(w).SetWriteDeadline(
				transferDeadline(app.Write.Duration,
					int64(file.Bytes)))
		}

		file.Serve(w)
		app.Log.Info("served file",
			"id", file.ID,
			"name", file.Name,
			"size", file.Size,
			"downloads", file.Count,
			"user", req)

		if reason != "" {
			app.Log.Info("removed file",
				"reason", reason,
				"id", file.ID,
				"name", file.Name,
				"downloads", file.Count)
		}
	}
}
