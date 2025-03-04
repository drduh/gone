package v1

import (
	"os"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
	"github.com/drduh/gone/signal"
)

func Run() {
	app := config.Load()

	app.Log.Info("started v1",
		"version", app.Version, "host", app.Hostname)

	signal.Setup(app)

	if err := server.Serve(app); err != nil {
		app.Log.Error("failed to start server",
			"error", err.Error())
		os.Exit(1)
	}
}
