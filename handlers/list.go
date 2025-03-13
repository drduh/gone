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
			file := config.File{
				Name:     record.Name,
				Size:     record.Size,
				Uploaded: record.Uploaded,
				Owner: config.Owner{
					Address: record.Owner.Address,
					Agent:   record.Owner.Agent,
				},
			}
			files = append(files, file)
		}

		writeJSON(w, http.StatusOK, files)
		app.Log.Info("served file list",
			"files", len(files),
			"ip", ip, "ua", ua)
	}
}
