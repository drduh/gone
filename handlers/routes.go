package handlers

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Routes returns the route table to register mux handlers.
func Routes(app *config.App) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{

		// Index page
		app.Root: Index(app),

		// Service status
		app.Status: Status(app),

		// Storage clear
		app.Clear: Clear(app),

		// Files
		app.Download: Download(app),
		app.List:     List(app),
		app.Upload:   Upload(app),

		// Messages
		app.Message:      Message(app),
		app.MessageClear: MessageClear(app),

		// Wall
		app.Wall: Wall(app),

		// User request information
		app.UserInfo:   UserInfo(app),
		app.UserRemask: UserRemask(app),

		// Misc
		app.Random: Random(app),
		app.Static: Static(app),
	}
}
