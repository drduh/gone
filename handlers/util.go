package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// JSON response helper for deny (auth fail)
func deny(w http.ResponseWriter, app *config.App, r *Request) {
	writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
	app.Log.Error(errorDeny, "user", r)
}

// Writes JSON response
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// Writes File response with content
func writeFile(w http.ResponseWriter, f *config.File) {
	w.Header().Set("Content-Type", f.MimeType())
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(f.Data)
}

// Returns CSS theme based on current time of day if unset
func getTheme(theme string) string {
	if theme != "" {
		return theme
	}
	if util.IsDaytime() {
		return "light"
	}
	return "dark"
}
