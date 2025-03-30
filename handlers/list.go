package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Returns list of file records
func List(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if !isAllowed(app, r) {
			deny(w, app, req)
			return
		}

		if !app.Throttle.Allow(app.Settings.Limits.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.Error.RateLimit))
			app.Log.Error(app.Error.RateLimit, "user", req)
			return
		}

		files := make([]config.File, 0, len(app.Storage.Files))
		for _, file := range app.Storage.Files {
			reason := file.IsExpired(app.Settings)
			if reason != "" {
				app.Storage.Expire(file)
				app.Log.Info("removed file",
					"reason", reason, "filename", file.Name,
					"downloads", file.Downloads.Total)
			} else {
				f := config.File{
					Name: file.Name,
					Size: file.Size,
					Owner: config.Owner{
						Address: file.Owner.Address,
						Agent:   file.Owner.Agent,
					},
					Time: config.Time{
						Upload: file.Upload,
						Remain: file.TimeRemaining().String(),
					},
					Downloads: config.Downloads{
						Allow:  file.Downloads.Allow,
						Total:  file.Downloads.Total,
						Remain: file.NumRemaining(),
					},
				}
				files = append(files, f)
			}
		}

		app.Log.Info("serving file list",
			"files", len(files), "user", req)
		writeJSON(w, http.StatusOK, files)
	}
}
