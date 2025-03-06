package auth

import "net/http"

const HEADER = "X-Auth"

// Returns true if request header contains valid pass
func Basic(pass string, r *http.Request) bool {

	// Always allow access if pass is not configured
	if pass == "" {
		return true
	}

	// Check header for non-empty pass and validate
	passHeader := r.Header.Get(HEADER)
	if passHeader != "" {
		return passHeader == pass
	}

	return false
}
