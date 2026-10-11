package auth

import (
	"net/http"
	"time"
)

// GetCookie returns a cookie value, creating one with
// the defaultValue, if none exists.
func GetCookie(w http.ResponseWriter, r *http.Request,
	defaultValue, id string, t time.Duration) string {
	cookie, err := r.Cookie(id)
	if err != nil || cookie.Value == "" {
		http.SetCookie(w, NewCookie(
			defaultValue, id, t, IsTLS(r)))
		return defaultValue
	}

	return cookie.Value
}

// IsTLS returns true if the request arrived over TLS.
func IsTLS(r *http.Request) bool {
	return r.TLS != nil
}

// NewCookie returns an HTTP cookie with the requested
// id, value and duration/expiration.
// Secure is only set over TLS.
func NewCookie(value, id string, t time.Duration,
	secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     id,
		Value:    value,
		Path:     "/",
		MaxAge:   int(t.Seconds()),
		Expires:  time.Now().Add(t),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
}
