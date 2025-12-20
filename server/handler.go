package server

import (
	"net/http"
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

const assetsPath = "assets"

// getHandler configures paths to handle and serve.
func getHandler(app *config.App) http.Handler {
	mux := http.NewServeMux()

	handle := func(path string, h http.HandlerFunc) {
		if path != "" {
			mux.HandleFunc(path, h)
		}
	}

	if app.Assets != "" {
		if _, err := os.Stat(assetsPath); err == nil {
			mux.Handle(app.Assets,
				http.StripPrefix(app.Assets,
					http.FileServer(http.Dir(assetsPath))))
		}
	}

	handle(app.Clear, handlers.Clear(app))
	handle(app.Download, handlers.Download(app))
	handle(app.List, handlers.List(app))
	handle(app.Message, handlers.Message(app))
	handle(app.Random, handlers.Random(app))
	handle(app.Root, handlers.Index(app))
	handle(app.Static, handlers.Static(app))
	handle(app.Status, handlers.Status(app))
	handle(app.Upload, handlers.Upload(app))
	handle(app.Wall, handlers.Wall(app))

	return mux
}
