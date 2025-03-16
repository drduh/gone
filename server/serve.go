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

	if app.Settings.Paths.Heartbeat != "" {
		mux.HandleFunc(app.Settings.Paths.Heartbeat, handlers.Heartbeat(app))
	}

	if app.Settings.Paths.Download != "" {
		mux.HandleFunc(app.Settings.Paths.Download, handlers.Download(app))
	}

	if app.Settings.Paths.List != "" {
		mux.HandleFunc(app.Settings.Paths.List, handlers.List(app))
	}

	if app.Settings.Paths.Static != "" {
		mux.HandleFunc(app.Settings.Paths.Static, handlers.Static(app))
	}

	if app.Settings.Paths.Upload != "" {
		mux.HandleFunc(app.Settings.Paths.Upload, handlers.Upload(app))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	app.Log.Info("starting server",
		"port", app.Port)

	return srv.ListenAndServe()
}
