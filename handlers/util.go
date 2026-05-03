package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

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
		app.List:     app.Require.List,
		app.Message:  app.Require.Message,
		app.Random:   app.Require.Random,
		app.Root:     app.Require.Root,
		app.Static:   app.Require.Static,
		app.Status:   app.Require.Status,
		app.Upload:   app.Require.Upload,
		app.User:     app.Require.User,
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
		IsBrowser: strings.Contains(
			r.Header.Get("Accept"), "text/html",
		),
	}
}

// authRequest returns only allowed parsed http Requests,
// rejecting unauthenticated and unauthorized attempts.
func authRequest(w http.ResponseWriter,
	r *http.Request, app *config.App) *Request {
	req := parseRequest(r)
	if !isAuthenticated(app, r) {
		deny(w, http.StatusForbidden, app.Deny)
		app.Log.Error(app.Deny, "user", req)
		return nil
	}
	if !app.Authorize(app.ReqsPerMinute) {
		deny(w, http.StatusTooManyRequests, app.RateLimit)
		app.Log.Error(app.RateLimit, "user", req)
		return nil
	}
	return req
}
