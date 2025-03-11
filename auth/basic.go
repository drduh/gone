package auth

import (
	"fmt"
	"net/http"
)

// Returns true if request header contains valid token
func Basic(header, token string, r *http.Request) bool {
	fmt.Println("header: ", header)

	// Always allow access if token is not configured
	if token == "" {
		return true
	}

	fmt.Println("header: ", header)

	// Check header for non-empty token and validate
	tokenHeader := r.Header.Get(header)
	fmt.Println("got token: ", tokenHeader)
	if tokenHeader != "" {
		fmt.Println("got token: ", tokenHeader)
		return tokenHeader == token
	}

	// Check form field value
	tokenForm := r.FormValue(header)
	if tokenForm != "" {
		fmt.Println("got token: ", tokenHeader)
		return tokenForm == token
	}

	fmt.Println("no token, returning false")

	return false
}
