package config

import (
	"log"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/util"
	"github.com/drduh/gone/version"
)

// Load returns the application configuration.
func Load() *App {
	app := App{}

	app.Start()
	app.Debug = modeDebug
	app.Modes.Version = modeVersion
	app.Settings = settings.Load(pathConfig)

	auditor, err := audit.Start(&audit.Config{
		Debug:      app.Debug || app.LogDebug,
		Filename:   app.LogFilename,
		TimeFormat: app.TimeFormat,
	})
	if err != nil {
		log.Fatalf("failed starting auditor: %v", err)
	}
	app.Log = auditor.Log

	app.Hostname = util.GetHostname()
	app.Version = version.Get()
	app.ClearStorage()

	auth.SetTarpit(app.TarpitDelay.GetDuration())

	return &app
}
