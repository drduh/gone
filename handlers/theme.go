package handlers

import (
	"net/http"
	"slices"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/util"
)

const autoTheme = "auto"

// getDefaultTheme returns a default theme, based on
// the current time if set to automatically theme.
func getDefaultTheme(theme string) string {
	if theme != autoTheme {
		return theme
	}
	if util.IsDaytime() {
		return "light"
	}
	return "dark"
}

// getTheme returns the CSS theme based on cookie preference,
// setting the cookie value if none exists, or is invalid.
func getTheme(w http.ResponseWriter, r *http.Request,
	defaultTheme, id string, t time.Duration, themes []string) string {
	formContent := r.FormValue(formFieldTheme)
	if formContent != "" {
		theme := formContent
		if !slices.Contains(themes, theme) {
			theme = getDefaultTheme(autoTheme)
		}
		http.SetCookie(w, auth.NewCookie(theme, id, t))
		return theme
	}
	return auth.GetCookie(w, r, defaultTheme, id, t)
}
