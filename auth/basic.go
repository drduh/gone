package auth

import "net/http"

// Basic returns true if an HTTP request header contains
// the valid header and token combination.
func Basic(header, token string, r *http.Request) bool {

	// Always allow access if token is not configured
	if token == "" {
		return true
	}

	// Check header for non-empty token and validate
	tokenHeader := r.Header.Get(header)
	if tokenHeader != "" {
		return tokenHeader == token
	}

	// Check form field value
	tokenForm := r.FormValue(header)
	if tokenForm != "" {
		return tokenForm == token
	}

	return false
}
