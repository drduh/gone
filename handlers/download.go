package handlers

import (
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// Returns content by file name
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if app.Settings.Auth.Require.Download &&
			!auth.Basic(app.Settings.Auth.Basic, r) {
			writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
			app.Log.Error(errorDeny,
				"action", "download",
				"ip", ip, "ua", ua)
			return
		}

		fileName := r.URL.Query().Get("name")

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
				"file", fileName,
				"ip", ip, "ua", ua)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(file.Data)

		file.Downloads++
		app.Log.Info("file downloaded",
			"filename", file.Name,
			"size", file.Size,
			"downloads", file.Downloads,
			"ip", ip, "ua", ua)

		expireReason := file.IsExpired(app.Settings)
		if expireReason != "" {
			delete(app.Storage.Files, file.Name)
			app.Log.Info("removed file",
				"reason", expireReason,
				"filename", file.Name,
				"downloads", file.Downloads)
		}
	}
}
