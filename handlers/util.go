package handlers

import (
	"net/http"
	"strings"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// authRequest returns only allowed parsed http Requests,
// rejecting unauthenticated and unauthorized attempts.
func authRequest(w http.ResponseWriter,
	r *http.Request, app *config.App) *Request {
	req := parseRequest(r)

	if !isAuthenticated(app, r) {
		app.Log.Error(app.Deny,
			"user", req)
		auth.ApplyTarpit()
		deny(w, http.StatusForbidden, app.Deny)
		return nil
	}

	if !app.Authorize(app.ReqsPerMinute) {
		app.Log.Error(app.RateLimit,
			"limit", app.ReqsPerMinute,
			"user", req)
		auth.ApplyTarpit()
		deny(w, http.StatusTooManyRequests, app.RateLimit)
		return nil
	}

	return req
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
		app.UserInfo: app.Require.UserInfo,
		app.Wall:     app.Require.Wall,
	}

	path := util.GetBasePath(r.URL.Path)
	required, exists := reqs[path]
	app.Log.Debug("checking authn",
		"path", path,
		"required", required,
		"exists", exists)

	if !exists || !required {
		app.Log.Debug("authn not required",
			"path", path)
		return true
	}

	token := getToken(app.Basic.Field, r)
	secret := []byte(app.Basic.Token)
	return auth.Basic(secret, token)
}

// parseRequest returns a Request with masked address.
func parseRequest(r *http.Request) *Request {
	mask := getMask(r.RemoteAddr)
	return &Request{
		Address: r.RemoteAddr,
		Agent:   r.UserAgent(),
		Mask:    mask,
		Path:    r.URL.String(),
		IsBrowser: strings.Contains(
			r.Header.Get("Accept"), "text/html",
		),
	}
}

// toPath redirects an HTTP request to a path.
func toPath(w http.ResponseWriter,
	r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// getMask returns a masked address string.
func getMask(addr string) string {
	return util.GetMaskAddr(addr, false)
}

// refreshMask sets a new masked address.
func refreshMask(addr string) {
	util.GetMaskAddr(addr, true)
}
