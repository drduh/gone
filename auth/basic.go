package auth

import "net/http"

// Basic returns true if request header or form value
// matches configured token.
func Basic(header, token string, r *http.Request) bool {

	// Allow access if token is not configured
	if token == "" {
		return true
	}

	// Check header for non-empty token and validate
	tokenHeader := r.Header.Get(header)
	if tokenHeader != "" && tokenHeader == token {
		return true
	}

	// Check form field value
	tokenForm := r.FormValue(header)
	if tokenForm != "" && tokenForm == token {
		return true
	}

	// Slow failed attempts
	applyTarpit()

	return false
}
