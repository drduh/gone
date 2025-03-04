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

		files := make([]templates.List, 0, len(app.Storage.Files))
		for _, rec := range app.Storage.Files {
			file := templates.List{
				Name: rec.Name,
				Size: rec.Size,
				Owner: templates.Owner{
					Address:  rec.Owner.Address,
					Agent:    rec.Owner.Agent,
					Uploaded: rec.Uploaded,
				},
			}
			files = append(files, file)
		}

		writeJSON(w, http.StatusOK, files)
		app.Log.Info("listed files",
			"files", len(files), "ip", ip, "ua", ua)
	}
}
