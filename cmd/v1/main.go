package v1

import (
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
)

func Run() {
	app := config.Load()

	app.Log.Info("started server",
		"version", app.Version,
		"host", app.Hostname,
		"port", app.Settings.Port)

	if err := server.Serve(app); err != nil {
		app.Log.Error("failed to start server",
			"error", err.Error())
		os.Exit(1)
	}
}
