package handlers

import (
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// AuthRequest returns an allowed, parsed HTTP Request,
// rejecting unauthenticated and unauthorized attempts.
func AuthRequest(
	w http.ResponseWriter,
	r *http.Request,
	app *config.App,
) *Request {
	req := parseRequest(r)

	if !checkAuthenticated(w, r, app) {
		app.Log.Error(app.Deny,
			"field", app.Basic.Field,
			"user", req)
		return nil
	}

	if !checkAuthorized(w, app) {
		app.Log.Error(app.RateLimit,
			"limit", app.ReqsPerMinute,
			"user", req)
		return nil
	}

	return req
}

func checkAuthenticated(
	w http.ResponseWriter,
	r *http.Request,
	app *config.App) bool {
	if !isAuthenticated(app, r) {
		auth.ApplyTarpit()
		deny(w, http.StatusForbidden, app.Deny)
		return false
	}
	return true
}

func checkAuthorized(
	w http.ResponseWriter,
	app *config.App) bool {
	if !app.Authorize(app.ReqsPerMinute) {
		auth.ApplyTarpit()
		deny(w, http.StatusTooManyRequests, app.RateLimit)
		return false
	}
	return true
}

// getToken reads the token from the request header,
// or posted form value.
func getToken(field string, r *http.Request) []byte {
	value := r.Header.Get(field)

	if value == "" {
		value = r.PostFormValue(field)
	}

	return []byte(value)
}

// isAuthenticated returns true if authentication for a route
// is configured and required.
func isAuthenticated(app *config.App, r *http.Request) bool {
	reqs := map[string]bool{
		app.Assets:       app.Require.Assets,
		app.Clear:        app.Require.Clear,
		app.Download:     app.Require.Download,
		app.List:         app.Require.List,
		app.Message:      app.Require.Message,
		app.MessageClear: app.Require.MessageClear,
		app.Random:       app.Require.Random,
		app.Root:         app.Require.Root,
		app.Static:       app.Require.Static,
		app.Status:       app.Require.Status,
		app.Upload:       app.Require.Upload,
		app.UserInfo:     app.Require.UserInfo,
		app.UserRemask:   app.Require.UserRemask,
		app.Wall:         app.Require.Wall,
	}

	path := r.Pattern
	required, exists := reqs[path]

	if !exists {
		app.Log.Debug("deny - no auth policy",
			"path", path)
		return false
	}

	if !required {
		app.Log.Debug("pass - auth not required",
			"path", path)
		return true
	}

	field := getToken(app.Basic.Field, r)
	secret := []byte(app.Basic.Token)

	app.Log.Debug("checking basic auth",
		"path", path,
		"required", required,
		"exists", exists)
	return auth.Basic(secret, field)
}
