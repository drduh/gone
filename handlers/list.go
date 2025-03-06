package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Returns list of file records
func List(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		files := make([]templates.File, 0, len(app.Storage.Files))
		for _, record := range app.Storage.Files {
			file := templates.File{
				Name: record.Name,
				Size: record.Size,
				Owner: templates.Owner{
					Address:  record.Owner.Address,
					Agent:    record.Owner.Agent,
					Uploaded: record.Uploaded,
				},
			}
			files = append(files, file)
		}

		writeJSON(w, http.StatusOK, files)
		app.Log.Info("files listed",
			"files", len(files),
			"ip", ip, "ua", ua)
	}
}
