package handlers

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// Writes JSON-encoded response
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// Writes File response with content
func writeFile(w http.ResponseWriter, data []byte, name string) {
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(data)
}

// Returns CSS theme based on current time of day if unset
func getTheme(theme string) string {
	if theme != "" {
		return theme
	}
	now := time.Now().Hour()
	if now > 7 && now < 19 {
		return "light"
	}
	return "dark"
}
