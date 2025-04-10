package auth

import (
	"net/http"
	"time"
)

// GetCookie returns a cookie value,
// creating one with the defaultValue if none exists.
func GetCookie(w http.ResponseWriter, r *http.Request,
	defaultValue, id string, t time.Duration) string {
	cookie, err := r.Cookie(id)
	if err != nil || cookie.Value == "" {
		http.SetCookie(w, NewCookie(defaultValue, id, t))
		return defaultValue
	}
	return cookie.Value
}

// NewCookie returns an HTTP cookie with
// a requested id, value and expiration.
func NewCookie(value, id string, t time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:    id,
		Expires: time.Now().Add(t),
		Path:    "/",
		Value:   value,
	}
}
