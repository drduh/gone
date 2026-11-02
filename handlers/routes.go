package handlers

import (
	"net/http"
	"strings"

	"github.com/drduh/gone/config"
)

// post prefixes a pattern with the POST method.
func post(pattern string) string {
	if pattern == "" {
		return ""
	}

	return http.MethodPost + " " + pattern
}

// patternPath returns the path of a mux pattern, without
// method prefix or "{$}" suffix.
func patternPath(pattern string) string {
	if _, path, ok := strings.Cut(pattern, " "); ok {
		pattern = path
	}

	return strings.TrimSuffix(pattern, "{$}")
}

// Register adds enabled routes to mux.
func Register(mux *http.ServeMux, app *config.App) {
	for pattern, h := range Routes(app) {
		if pattern == "" {
			continue
		}

		mux.HandleFunc(pattern, h)
	}
}

// Routes returns the route table to register mux handlers.
func Routes(app *config.App) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{

		// Index page (exact match)
		app.Root + "{$}": Index(app),

		// Service status
		app.Status: Status(app),

		// Storage clear
		post(app.Clear): Clear(app),

		// Files
		app.Download:     Download(app),
		app.List:         List(app),
		post(app.Upload): Upload(app),

		// Messages
		app.Message:             Message(app),
		post(app.MessageAdd):    MessageAdd(app),
		post(app.MessageClear):  MessageClear(app),
		post(app.MessageDelete): MessageDelete(app),
		app.MessageGet:          MessageGet(app),

		// Wall
		app.WallModify: Wall(app),

		// User request information
		app.UserInfo:         UserInfo(app),
		post(app.UserRemask): UserRemask(app),

		// Random string generator
		app.Random: Random(app),

		// Static embedded content
		app.Static: Static(app),

		// Remaining undefined routes
		"/": Undefined(app),
	}
}
