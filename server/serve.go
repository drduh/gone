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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	app.Log.Info("starting server", "port", app.Port)

	return srv.ListenAndServe()
}
