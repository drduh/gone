package server

import (
	"fmt"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// Start HTTP server with configured paths
func Serve(app *config.App) error {
	go expiryWorker(app)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Index(app))

	paths := app.Settings.Paths
	if paths.Heartbeat != "" {
		mux.HandleFunc(paths.Heartbeat, handlers.Heartbeat(app))
	}
	if paths.Download != "" {
		mux.HandleFunc(paths.Download, handlers.Download(app))
	}
	if paths.List != "" {
		mux.HandleFunc(paths.List, handlers.List(app))
	}
	if paths.Static != "" {
		mux.HandleFunc(paths.Static, handlers.Static(app))
	}
	if paths.Upload != "" {
		mux.HandleFunc(paths.Upload, handlers.Upload(app))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	app.Log.Info("starting server",
		"port", app.Port)

	return srv.ListenAndServe()
}
