package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// List handles requests to list Files in Storage.
func List(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if !app.Allow(app.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.RateLimit))
			app.Log.Error(app.RateLimit, "user", req)
			return
		}

		app.UpdateTime()
		files := getFiles(app)

		app.Log.Info("serving file list",
			"files", len(files), "user", req)
		writeJSON(w, http.StatusOK, files)
	}
}

// getFiles returns a list of non-expired Files in Storage.
func getFiles(app *config.App) []storage.File {
	files := make([]storage.File, 0, len(app.Files))
	for _, file := range app.Files {
		reason := file.IsExpired()
		if reason != "" {
			app.Expire(file)
			app.Log.Info("removed file",
				"reason", reason, "filename", file.Name,
				"downloads", file.Total)
			break
		}
		f := storage.File{
			Name: file.Name,
			Size: file.Size,
			Type: file.Type,
			Owner: storage.Owner{
				Address: file.Address,
				Agent:   file.Agent,
			},
			Time: storage.Time{
				Remain: file.Time.Remain,
				Upload: file.Upload,
			},
			Downloads: storage.Downloads{
				Allow:  file.Downloads.Allow,
				Remain: file.NumRemaining(),
				Total:  file.Total,
			},
		}
		files = append(files, f)
	}
	return files
}
