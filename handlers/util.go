package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// deny serves a JSON response for deny (auth fail).
func deny(w http.ResponseWriter, app *config.App, r *Request) {
	writeJSON(w, http.StatusUnauthorized, errorJSON(app.Deny))
	app.Log.Error(app.Deny, "user", r)
}

// toRoot redirects an HTTP request to the root/index handler.
func toRoot(w http.ResponseWriter, r *http.Request, rootPath string) {
	http.Redirect(w, r, rootPath, http.StatusSeeOther)
}

// isAllowed returns true if authentication for a route
// is required and allowed.
func isAllowed(app *config.App, r *http.Request) bool {
	reqs := map[string]bool{
		app.Download: app.Require.Download,
		app.List:     app.Require.List,
		app.Message:  app.Require.Message,
		app.Upload:   app.Require.Upload,
		app.Wall:     app.Require.Wall,
	}

	path := util.GetBasePath(r.URL.Path)

	required, exists := reqs[path]
	app.Log.Debug("checking auth",
		"path", path, "required", required, "exists", exists)
	if !exists || !required {
		app.Log.Debug("auth not required", "path", r.URL.Path)
		return true
	}

	return isAuthenticated(app.Basic.Field, app.Basic.Token, r)
}

// isAuthentication returns true if basic authentication is successful.
func isAuthenticated(header, token string, r *http.Request) bool {
	return auth.Basic(header, token, r)
}

// errorJSON returns an error string map containing the string.
func errorJSON(s string) map[string]string {
	return map[string]string{
		"error": s,
	}
}

// getDefaultTheme returns a theme based on
// the current time of day if set to "auto".
func getDefaultTheme(theme string) string {
	if theme != "auto" {
		return theme
	}
	if util.IsDaytime() {
		return "light"
	}
	return "dark"
}

// parseRequest returns a Request with masked address.
func parseRequest(r *http.Request) *Request {
	address, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		address = "unknown address"
	}
	return &Request{
		Action:  r.URL.String(),
		Address: r.RemoteAddr,
		Mask:    util.Mask(address),
		Agent:   r.UserAgent(),
	}
}

// authRequest returns only an allowed parsed http Request struct.
func authRequest(w http.ResponseWriter, r *http.Request, app *config.App) (*Request, bool) {
	req := parseRequest(r)
	if !isAllowed(app, r) {
		deny(w, app, req)
		return nil, false
	}
	return req, true
}

// writeJSON serves a JSON response with data.
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

// getParam return a parameter (such as filename)
// from an HTTP request, as URL, query or form value.
func getParam(r *http.Request, pathLen int, fieldName string) string {
	p := r.URL.Path[pathLen:]
	if p == "" {
		p = r.URL.Query().Get(fieldName)
	}
	if p == "" {
		p = r.FormValue(fieldName)
	}
	return p
}

// getTheme returns the CSS theme based on cookie preference,
// setting the cookie value if none exists.
func getTheme(w http.ResponseWriter, r *http.Request,
	defaultTheme, id string, t time.Duration) string {
	theme := r.FormValue("theme")
	if theme != "" {
		http.SetCookie(w, auth.NewCookie(theme, id, t))
		return theme
	}
	return auth.GetCookie(w, r, defaultTheme, id, t)
}
