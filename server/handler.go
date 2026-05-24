package server

import (
	"net/http"
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// getHandler configures paths to handle and serve.
func getHandler(app *config.App) http.Handler {
	const pathAssets = "assets"

	mux := http.NewServeMux()

	handle := func(path string, h http.HandlerFunc) {
		if path != "" {
			mux.HandleFunc(path, h)
		}
	}

	if app.Assets != "" {
		if _, err := os.Stat(pathAssets); err == nil {
			app.Log.Debug("assets present", "path", pathAssets)
			mux.Handle(app.Assets, http.StripPrefix(
				app.Assets, http.FileServer(http.Dir(pathAssets))))
		} else {
			app.Log.Warn("missing assets", "path", pathAssets)
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
	handle(app.User, handlers.User(app))
	handle(app.Wall, handlers.Wall(app))

	return mux
}
