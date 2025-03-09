package server

import (
	"fmt"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// Start HTTP server with application configuration
func Serve(app *config.App) error {
	go expiryWorker(app)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Heartbeat(app))
	mux.HandleFunc(app.Settings.Paths.Download, handlers.Download(app))
	mux.HandleFunc(app.Settings.Paths.List, handlers.List(app))
	mux.HandleFunc(app.Settings.Paths.Static, handlers.Static(app))
	mux.HandleFunc(app.Settings.Paths.Upload, handlers.Upload(app))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	app.Log.Info("starting server",
		"port", app.Port)

	return srv.ListenAndServe()
}
