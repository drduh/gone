// Package v1 implements the original design circa 2025.
package v1

import (
	"flag"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
	"github.com/drduh/gone/signal"
	"github.com/drduh/gone/version"
)

// Run loads the application configuration, sets up the
// signal handler and starts server, or prints version.
func Run() int {
	flag.Parse()

	app := config.Load()

	if app.Modes.Version {
		version.Print()
		return 0
	}

	app.Log.Info("starting v1",
		"host", app.Hostname,
		"version", app.Version)
	app.Log.Debug("debug log enabled",
		"configuration", app)

	signal.Setup(app)

	if err := server.Serve(app); err != nil {
		app.Log.Error("failed to start v1",
			"error", err.Error())
		return 1
	}

	return 0
}
