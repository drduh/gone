package server

import (
	"net/http"
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// getHandler configures HTTP handler routes to serve.
func getHandler(app *config.App) http.Handler {
	mux := http.NewServeMux()

	handle := func(path string, h http.HandlerFunc) {
		if path != "" {
			mux.HandleFunc(path, h)
		}
	}

	handle("/", handlers.Index(app))

	if app.Assets != "" {
		assets := "assets"
		if _, err := os.Stat(assets); err == nil {
			mux.Handle(app.Assets, http.StripPrefix(
				app.Assets, http.FileServer(http.Dir(assets))))
		}
	}

	handle(app.Download, handlers.Download(app))
	handle(app.Heartbeat, handlers.Heartbeat(app))
	handle(app.List, handlers.List(app))
	handle(app.Message, handlers.Message(app))
	handle(app.Static, handlers.Static(app))
	handle(app.Upload, handlers.Upload(app))

	return mux
}
