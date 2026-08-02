package handlers

import (
	"net/http"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// Undefined applies tarpit to requests to undefined paths.
func Undefined(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.ApplyTarpit()

		req := parseRequest(r)

		app.Log.Info(app.NoPath,
			"method", r.Method,
			"user", req)
		deny(w, http.StatusNotFound, app.NoPath)
	}
}
