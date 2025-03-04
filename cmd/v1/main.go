package v1

import (
	"time"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/config"
)

func Run() {
	app := config.Load()

	l, _ := audit.StartAuditor()

	l.Log.Info("started server",
		"version", app.Version,
		"host", app.Hostname,
		"runtime", time.Since(app.Start).String())
}
