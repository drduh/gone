package auth

import "net/http"

const HEADER = "X-Auth"

// Returns true if request header contains valid token
func Basic(token string, r *http.Request) bool {

	// Always allow access if token is not configured
	if token == "" {
		return true
	}

	// Check header for non-empty token and validate
	tokenHeader := r.Header.Get(HEADER)
	if tokenHeader != "" {
		return tokenHeader == token
	}

	return false
}
