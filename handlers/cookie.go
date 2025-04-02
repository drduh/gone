package handlers

import (
	"net/http"
	"time"
)

// Returns HTTP cookie with requested id and expiration
func newCookie(id string, t time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:    id,
		Expires: time.Now().Add(t),
		Path:    "/",
	}
}
