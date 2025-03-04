package v1

import (
	"os"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/server"
)

func Run() {
	app := config.Load()

	app.Log.Info("started server",
		"version", app.Version,
		"host", app.Hostname,
		"runtime", time.Since(app.Start).String())

	if err := server.Serve(); err != nil {
		app.Log.Error("failed to start server",
			"error", err.Error())
		os.Exit(1)
	}
}
