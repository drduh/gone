package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// getHandler configures paths to handle and serve.
func getHandler(app *config.App) http.Handler {
	mux := http.NewServeMux()

	registerAssets(mux, app)
	handlers.Register(mux, app)

	return withCSRF(app, mux)
}

// withCSRF rejects cross-origin browser form posts
func withCSRF(app *config.App, h http.Handler) http.Handler {
	csrf := http.NewCrossOriginProtection()
	csrf.SetDenyHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.Log.Error(app.Deny,
				"reason", "cross-origin",
				"origin", r.Header.Get("Origin"),
				"path", r.URL.Path)
			handlers.Deny(w, http.StatusForbidden, app.Deny)
		}))

	return csrf.Handler(h)
}

// registerAssets sets up HTML asset handling.
func registerAssets(mux *http.ServeMux, app *config.App) {
	const pathAssets = "assets"

	if app.Assets == "" {
		return
	}

	if _, err := os.Stat(pathAssets); err != nil {
		app.Log.Warn("missing assets",
			"path", pathAssets)
		return
	}

	assetsPath := app.Assets
	if !strings.HasSuffix(assetsPath, "/") {
		assetsPath += "/"
	}

	mux.Handle(assetsPath, wrapAssets(
		app,
		http.StripPrefix(
			assetsPath,
			http.FileServer(http.Dir(pathAssets))),
	))
}

// wrapAssets applies rate-limiting and caching to assets.
func wrapAssets(app *config.App, h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			req := handlers.AuthRequest(w, r, app)
			if req == nil {
				return
			}

			w.Header().Set("Cache-Control", "no-cache")
			h.ServeHTTP(w, r)
		})
}
