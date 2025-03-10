package v1

import (
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
	"github.com/drduh/gone/signal"
	"github.com/drduh/gone/version"
)

func Run() {
	app := config.Load()

	if app.Modes.Version {
		version.Print()
		os.Exit(0)
	}

	app.Log.Info("started v1",
		"version", app.Version, "host", app.Hostname)
	app.Log.Debug("debug logging enabled",
		"configuration", app)

	signal.Setup(app)

	if err := server.Serve(app); err != nil {
		app.Log.Error("failed to start server",
			"error", err.Error())
		os.Exit(1)
	}
}
