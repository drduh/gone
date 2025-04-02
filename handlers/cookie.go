package handlers

import (
	"net/http"
	"time"
)

// Returns Theme based on user cookie preference
func getTheme(w http.ResponseWriter, r *http.Request, defaultTheme, id string, t time.Duration) string {
	theme := r.FormValue("theme")
	if theme != "" {
		http.SetCookie(w, newCookie(theme, id, t))
		return theme
	}
	return getCookie(w, r, defaultTheme, id, t)
}

// Gets or creates cookie value with theme
func getCookie(w http.ResponseWriter, r *http.Request, defaultTheme, id string, t time.Duration) string {
	cookie, err := r.Cookie(id)
	if err != nil || cookie.Value == "" {
		http.SetCookie(w, newCookie(defaultTheme, id, t))
		return defaultTheme
	}
	return cookie.Value
}

// Returns HTTP cookie with requested id and expiration
func newCookie(value, id string, t time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:    id,
		Expires: time.Now().Add(t),
		Path:    "/",
		Value:   value,
	}
}
