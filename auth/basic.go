package auth

import "net/http"

const HEADER = "X-Auth"

// Returns True if request header contains pass string
func Basic(pass string, r *http.Request) bool {

	// If pass is not configured, always return true
	if pass == "" {
		return true
	}

	// Check header for non-empty pass
	passHeader := r.Header.Get(HEADER)
	if passHeader != "" {
		return passHeader == pass
	}

	return false
}
