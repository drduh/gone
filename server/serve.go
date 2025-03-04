package server

import (
	"fmt"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

// Start HTTP server with application configuration
func Serve(app *config.App) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Heartbeat())
	mux.HandleFunc("/download", handlers.Download(app))
	mux.HandleFunc("/list", handlers.List(app))
	mux.HandleFunc("/static", handlers.Static(app))
	mux.HandleFunc("/upload", handlers.Upload(app))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	app.Log.Info("starting server", "port", app.Port)

	return srv.ListenAndServe()
}
