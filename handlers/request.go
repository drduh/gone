package handlers

import (
	"net/http"
	"strings"
)

// parseRequest returns a Request with masked address.
func parseRequest(r *http.Request) *Request {
	return &Request{
		Agent:       r.UserAgent(),
		Address:     r.RemoteAddr,
		AddressMask: getMask(r.RemoteAddr),
		Path:        r.URL.String(),
		IsBrowser: strings.Contains(
			r.Header.Get("Accept"), "text/html",
		),
	}
}

// toPath redirects an HTTP request to a path.
func toPath(w http.ResponseWriter, r *http.Request, p string) {
	http.Redirect(w, r, p, http.StatusSeeOther)
}
