// Package v1 implements the original (circa 2025) design.
package v1

import (
	"flag"
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
	"github.com/drduh/gone/signal"
	"github.com/drduh/gone/version"
)

// Run loads the application configuration, sets up the
// signal handler and starts the HTTP server.
func Run() {
	flag.Parse()

	app := config.Load()

	if app.Modes.Version {
		version.Print()
		os.Exit(0)
	}

	app.Log.Info("started v1",
		"version", app.Version, "host", app.Hostname)
	app.Log.Debug("debug enabled", "configuration", app)

	signal.Setup(app)

	if err := server.Serve(app); err != nil {
		app.Log.Error("server failed", "error", err.Error())
		os.Exit(1)
	}
}
