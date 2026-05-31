package auth

import "crypto/subtle"

// auth returns true if all bytes are equal.
func auth(secret, token []byte) bool {
	return subtle.ConstantTimeCompare(secret, token) == 1
}

// Basic returns true if token matches secret.
func Basic(secret, token []byte) bool {
	if len(token) == 0 {
		return false
	}

	if auth(secret, token) {
		return true
	}

	return false
}
