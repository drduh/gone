package server

import (
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// Configures HTTP handler routes
func getHandler(app *config.App) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index(app))
	p := app.Paths
	if p.Assets != "" {
		mux.Handle(p.Assets, http.StripPrefix(
			p.Assets, http.FileServer(http.Dir(
				"templates/data/assets/"))))
	}
	if p.Heartbeat != "" {
		mux.HandleFunc(p.Heartbeat, handlers.Heartbeat(app))
	}
	if p.Download != "" {
		mux.HandleFunc(p.Download, handlers.Download(app))
	}
	if p.List != "" {
		mux.HandleFunc(p.List, handlers.List(app))
	}
	if p.Message != "" {
		mux.HandleFunc(p.Message, handlers.Message(app))
	}
	if p.Static != "" {
		mux.HandleFunc(p.Static, handlers.Static(app))
	}
	if p.Upload != "" {
		mux.HandleFunc(p.Upload, handlers.Upload(app))
	}
	return mux
}
