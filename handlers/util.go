package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// deny serves a SON response for deny (auth fail)
func deny(w http.ResponseWriter, app *config.App, r *Request) {
	writeJSON(w, http.StatusUnauthorized, errorJSON(app.Deny))
	app.Log.Error(app.Deny, "user", r)
}

// isAllowed returns true if auth for route path is required and allowed
func isAllowed(app *config.App, r *http.Request) bool {
	reqs := map[string]bool{
		app.Download: app.Require.Download,
		app.List:     app.Require.List,
		app.Message:  app.Require.Message,
		app.Upload:   app.Require.Upload,
	}
	path := getBasePath(r.URL.Path)
	required, exists := reqs[path]
	app.Log.Debug("checking auth",
		"path", path, "required", required, "exists", exists)
	if !exists || !required {
		app.Log.Debug("auth not required", "path", r.URL.Path)
		return true
	}
	return isAuthenticated(app.Basic.Field, app.Basic.Token, r)
}

// isAuthentication returns true if an authentication is successful
func isAuthenticated(header, token string, r *http.Request) bool {
	return auth.Basic(header, token, r)
}

// errorJSON returns an error string map
func errorJSON(s string) map[string]string {
	return map[string]string{
		"error": s,
	}
}

// getDefaultTheme Returns a theme based on current time of day if set to "auto"
func getDefaultTheme(theme string) string {
	if theme != "auto" {
		return theme
	}
	if util.IsDaytime() {
		return "light"
	}
	return "dark"
}

// parseRequest returns parsed a HTTP Request struct for log
func parseRequest(r *http.Request) *Request {
	return &Request{
		Action:  r.URL.String(),
		Address: r.RemoteAddr,
		Agent:   r.UserAgent(),
	}
}

// authRequest returns an allowed parsed http Request struct
func authRequest(w http.ResponseWriter, r *http.Request, app *config.App) (*Request, bool) {
	req := parseRequest(r)
	if !isAllowed(app, r) {
		deny(w, app, req)
		return nil, false
	}
	return req, true
}

// writeJSON serves a JSON response with data
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

// getBasePath returns the string up to and including the first "/"
func getBasePath(s string) string {
	i := strings.Index(s[1:], "/")
	if i == -1 {
		return s
	}
	return s[:i+2]
}
