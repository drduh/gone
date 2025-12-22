package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// List handles requests to list Files in Storage.
func List(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		files := app.ListFiles()
		app.Log.Info("serving file list",
			"files", len(files), "user", req)
		writeJSON(w, http.StatusOK, files)
	}
}
