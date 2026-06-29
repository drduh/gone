package handlers

import (
	"net/http"
	"strings"

	"github.com/drduh/gone/util"
)

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
