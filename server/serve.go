package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// newServer sets up the HTTP server with
// configured timeouts and routes to handle.
func newServer(app *config.App) *http.Server {
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

	return &http.Server{
		Addr:              app.GetAddr(),
		Handler:           handler,
		IdleTimeout:       timeoutIdle,
		ReadHeaderTimeout: timeoutHeader,
		ReadTimeout:       timeoutRead,
		WriteTimeout:      timeoutWrite,
	}
}

// Serve starts the expiry worker and HTTP/S server.
func Serve(app *config.App) error {
	go expiryWorker(app)

	server := newServer(app)

	if (app.CertFile == "") != (app.KeyFile == "") {
		return errors.New(
			"TLS requires both certFile and keyFile",
		)
	}

	if app.CertFile != "" {
		app.Log.Info("starting HTTPS server",
			"addr", server.Addr,
			"certFile", app.CertFile,
			"keyFile", app.KeyFile)
		if err := server.ListenAndServeTLS(
			app.CertFile, app.KeyFile); err != nil {
			return fmt.Errorf("HTTPS server failed: %w", err)
		}
	} else {
		app.Log.Info("starting HTTP server",
			"addr", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			return fmt.Errorf("HTTP server failed: %w", err)
		}
	}

	return nil
}
