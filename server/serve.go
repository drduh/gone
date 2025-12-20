package server

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Serve starts the expiry worker and HTTP server
// with configured routes to handle.
func Serve(app *config.App) error {
	go expiryWorker(app)

	app.Log.Info("starting server", "port", app.Port)

	server := &http.Server{
		Addr:              app.GetAddr(),
		Handler:           getHandler(app),
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
	}

	return server.ListenAndServe()
}
