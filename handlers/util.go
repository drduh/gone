package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"slices"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

const autoTheme = "auto"

// deny serves a JSON response for disallowed requests.
func deny(w http.ResponseWriter, httpCode int, reason string,
	app *config.App, r *Request) {
	writeJSON(w, httpCode, errorJSON(reason))
	app.Log.Error(reason, "user", r)
}

// toRoot redirects an HTTP request to the root/index handler.
func toRoot(w http.ResponseWriter, r *http.Request, rootPath string) {
	http.Redirect(w, r, rootPath, http.StatusSeeOther)
}

// isAuthenticated returns true if authentication for a route
// is required and allowed.
func isAuthenticated(app *config.App, r *http.Request) bool {
	reqs := map[string]bool{
		app.Clear:    app.Require.Clear,
		app.Download: app.Require.Download,
		app.Root:     app.Require.Root,
		app.List:     app.Require.List,
		app.Message:  app.Require.Message,
		app.Upload:   app.Require.Upload,
		app.Wall:     app.Require.Wall,
	}

	path := util.GetBasePath(r.URL.Path)
	required, exists := reqs[path]
	app.Log.Debug("checking authn",
		"path", path, "required", required, "exists", exists)
	if !exists || !required {
		app.Log.Debug("authn not required", "path", r.URL.Path)
		return true
	}

	return auth.Basic(app.Basic.Field, app.Basic.Token, r)
}

// errorJSON returns an error string map containing the string.
func errorJSON(s string) map[string]string {
	return map[string]string{
		"error": s,
	}
}

// getDefaultTheme returns a default theme, based on
// the current time if set to automatically theme.
func getDefaultTheme(theme string) string {
	if theme != autoTheme {
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
		Address: r.RemoteAddr,
		Agent:   r.UserAgent(),
		Mask:    util.Mask(address),
		Path:    r.URL.String(),
	}
}

// authRequest returns only allowed parsed http Requests,
// rejecting unauthenticated and unauthorized attempts.
func authRequest(w http.ResponseWriter,
	r *http.Request, app *config.App) *Request {
	req := parseRequest(r)
	if !isAuthenticated(app, r) {
		deny(w, http.StatusForbidden, app.Deny, app, req)
		return nil
	}
	if !app.Authorize(app.ReqsPerMinute) {
		deny(w, http.StatusTooManyRequests, app.RateLimit, app, req)
		return nil
	}
	return req
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

// getRequestParameter returns a request parameter from the
// URL or a form value.
func getRequestParameter(r *http.Request,
	pathLen int, fieldName string) string {
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
// setting the cookie value if none exists, or is invalid.
func getTheme(w http.ResponseWriter, r *http.Request,
	defaultTheme, id string, t time.Duration, themes []string) string {
	formContent := r.FormValue(formFieldTheme)
	if formContent != "" {
		theme := formContent
		if !slices.Contains(themes, theme) {
			theme = getDefaultTheme(autoTheme)
		}
		http.SetCookie(w, auth.NewCookie(theme, id, t))
		return theme
	}
	return auth.GetCookie(w, r, defaultTheme, id, t)
}
