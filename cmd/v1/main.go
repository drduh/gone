package v1

import (
	"time"

	"github.com/drduh/gone/config"
)

func Run() {
	app := config.Load()

	app.Log.Info("started server",
		"version", app.Version,
		"host", app.Hostname,
		"runtime", time.Since(app.Start).String())
}
