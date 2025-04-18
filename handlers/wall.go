package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Wall handles wall form submissions and adds content to Storage.
func Wall(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			wallContent := r.FormValue("wall")
			if wallContent != "" {
				app.WallContent = wallContent
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
