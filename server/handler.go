package server

import (
	"net/http"
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// Configures HTTP handler routes
func getHandler(app *config.App) http.Handler {
	mux := http.NewServeMux()

	handle := func(path string, h http.HandlerFunc) {
		if path != "" {
			mux.HandleFunc(path, h)
		}
	}

	handle("/", handlers.Index(app))

	if app.Paths.Assets != "" {
		assets := "templates/data/assets/"
		if _, err := os.Stat(assets); err == nil {
			mux.Handle(app.Paths.Assets, http.StripPrefix(
				app.Paths.Assets, http.FileServer(http.Dir(assets))))
		}
	}

	handle(app.Paths.Download, handlers.Download(app))
	handle(app.Paths.Heartbeat, handlers.Heartbeat(app))
	handle(app.Paths.List, handlers.List(app))
	handle(app.Paths.Message, handlers.Message(app))
	handle(app.Paths.Static, handlers.Static(app))
	handle(app.Paths.Upload, handlers.Upload(app))

	return mux
}
