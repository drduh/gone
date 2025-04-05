package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// JSON response helper for deny (auth fail)
func deny(w http.ResponseWriter, app *config.App, r *Request) {
	writeJSON(w, http.StatusUnauthorized, errorJSON(app.Deny))
	app.Log.Error(app.Deny, "user", r)
}

// Returns true if auth for route path is required and allowed
func isAllowed(app *config.App, r *http.Request) bool {
	reqs := map[string]bool{
		app.Download: app.Require.Download,
		app.List:     app.Require.List,
		app.Message:  app.Require.Message,
		app.Upload:   app.Require.Upload,
	}
	app.Log.Debug("checking auth", "path", r.URL.Path)
	required, exists := reqs[r.URL.Path]
	if !exists || !required {
		return true
	}
	return isAuthenticated(app.Basic.Field, app.Basic.Token, r)
}

// Returns true if authentication is successful
func isAuthenticated(header, token string, r *http.Request) bool {
	return auth.Basic(header, token, r)
}

// Returns error in string map
func errorJSON(s string) map[string]string {
	return map[string]string{
		"error": s,
	}
}

// Returns CSS theme based on current time of day if unset
func getDefaultTheme(theme string) string {
	if theme != "" {
		return theme
	}
	if util.IsDaytime() {
		return "light"
	}
	return "dark"
}

// Returns parsed http Request struct for log
func parseRequest(r *http.Request) *Request {
	return &Request{
		Action:  r.URL.String(),
		Address: r.RemoteAddr,
		Agent:   r.UserAgent(),
	}
}

// Returns parsed http Request struct, if allowed
func authRequest(w http.ResponseWriter, r *http.Request, app *config.App) (*Request, bool) {
	req := parseRequest(r)
	if !isAllowed(app, r) {
		deny(w, app, req)
		return nil, false
	}
	return req, true
}

// Writes JSON response
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorJSON(err.Error()))
		return
	}
}
