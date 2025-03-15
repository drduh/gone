package handlers

import (
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// Returns list of file records
func List(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if app.Settings.Auth.Require.List &&
			!auth.Basic(app.Settings.Auth.Header, app.Settings.Auth.Token, r) {
			writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
			app.Log.Error(errorDeny,
				"action", "list",
				"ip", ip, "ua", ua)
			return
		}

		if throttle(app) {
			writeJSON(w, http.StatusTooManyRequests, responseErrorRateLimit)
			app.Log.Error(errorRateLimit,
				"action", "list",
				"ip", ip, "ua", ua)
			return
		}

		files := make([]config.File, 0, len(app.Storage.Files))
		for _, record := range app.Storage.Files {
			reason := record.IsExpired(app.Settings)
			if reason != "" {
				delete(app.Storage.Files, record.Name)
				app.Log.Info("removed file",
					"reason", reason,
					"filename", record.Name,
					"downloads", record.Downloads.Total)
			} else {
				file := config.File{
					Name: record.Name,
					Size: record.Size,
					Owner: config.Owner{
						Address: record.Owner.Address,
						Agent:   record.Owner.Agent,
					},
					Time: config.Time{
						Upload: record.Upload,
						Remain: record.TimeRemaining().String(),
					},
					Downloads: config.Downloads{
						Allow:  record.Downloads.Allow,
						Total:  record.Downloads.Total,
						Remain: record.NumRemaining(),
					},
				}
				files = append(files, file)
			}
		}

		writeJSON(w, http.StatusOK, files)
		app.Log.Info("served file list",
			"files", len(files),
			"ip", ip, "ua", ua)
	}
}
