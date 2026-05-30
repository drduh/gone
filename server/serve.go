package server

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Serve starts the expiry worker and HTTP server
// with timeouts and configured routes to handle.
func Serve(app *config.App) error {
	go expiryWorker(app)

	handler := getHandler(app)

	timeoutIdle := 90 * time.Second
	timeoutRead := 20 * time.Second
	timeoutHeader := 20 * time.Second
	timeoutWrite := 20 * time.Second
	app.Log.Debug("server timeouts",
		"idle", timeoutIdle.String(),
		"read", timeoutRead.String(),
		"header", timeoutHeader.String(),
		"write", timeoutWrite.String())

	address := app.GetAddr()
	app.Log.Info("starting server",
		"address", address)

	server := &http.Server{
		Addr:              address,
		Handler:           handler,
		IdleTimeout:       timeoutIdle,
		ReadHeaderTimeout: timeoutHeader,
		ReadTimeout:       timeoutRead,
		WriteTimeout:      timeoutWrite,
	}

	return server.ListenAndServe()
}
