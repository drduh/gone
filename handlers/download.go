package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Returns content by file name
func Download(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()
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
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "file not found"})
			app.Log.Error("file not found",
				"file", fileName, "ip", ip, "ua", ua)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(file.Data)

		app.Log.Info("served file",
			"name", file.Name, "size", file.Size,
			"ip", ip, "ua", ua)
	}
}
