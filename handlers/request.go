package handlers

import (
	"net"
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

// clientAddr returns the request host without port,
// or the raw RemoteAddr if it has no port.
func clientAddr(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

// toPath redirects an HTTP request to a path.
func toPath(w http.ResponseWriter, r *http.Request, p string) {
	http.Redirect(w, r, p, http.StatusSeeOther)
}
