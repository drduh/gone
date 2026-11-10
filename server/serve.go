package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/drduh/gone/config"
)

var errTLSMissingFiles = errors.New(
	"TLS requires both certificate and key files")

// newServer sets up the HTTP server with
// configured timeouts and routes to handle.
func newServer(app *config.App) *http.Server {
	handler := getHandler(app)

	timeoutIdle := app.Idle.Duration
	timeoutRead := app.Read.Duration
	timeoutHeader := app.ReadHeader.Duration
	timeoutWrite := app.Write.Duration
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

		ErrorLog: slog.NewLogLogger(
			app.Log.Handler(), slog.LevelError),
	}
}

// Serve starts the expiry worker and HTTP/S server.
func Serve(app *config.App) error {
	if (app.CertFile == "") != (app.KeyFile == "") {
		return errTLSMissingFiles
	}

	go expiryWorker(app)

	server := newServer(app)

	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	return serve(app, server, ln)
}

// serve runs server on ln until stopped via StopServer
func serve(app *config.App, server *http.Server, ln net.Listener) error {
	var draining sync.WaitGroup
	app.StopServer = func(ctx context.Context) error {
		draining.Add(1)
		defer draining.Done()
		return server.Shutdown(ctx)
	}

	var err error

	if app.CertFile != "" {
		server.TLSConfig = newTLSConfig()

		app.Log.Info("starting HTTPS server",
			"addr", ln.Addr().String(),
			"certFile", app.CertFile,
			"keyFile", app.KeyFile)
		err = server.ServeTLS(ln, app.CertFile, app.KeyFile)
	} else {
		app.Log.Info("starting HTTP server",
			"addr", ln.Addr().String())
		err = server.Serve(ln)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}

	draining.Wait()

	return nil
}
