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

		if !app.Allow(app.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.RateLimit))
			app.Log.Error(app.RateLimit, "user", req)
			return
		}

		files := make([]config.File, 0, len(app.Files))
		for _, file := range app.Files {
			reason := file.IsExpired(app.Settings)
			if reason != "" {
				app.Expire(file)
				app.Log.Info("removed file",
					"reason", reason, "filename", file.Name,
					"downloads", file.Total)
			} else {
				file.Time.Remain = file.TimeRemaining().String()
				app.Files[file.Name] = file
				f := config.File{
					Name: file.Name,
					Size: file.Size,
					Owner: config.Owner{
						Address: file.Address,
						Agent:   file.Agent,
					},
					Time: config.Time{
						Upload: file.Upload,
						Remain: file.Time.Remain,
					},
					Downloads: config.Downloads{
						Allow:  file.Downloads.Allow,
						Total:  file.Total,
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
