package config

import (
	"log"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/util"
	"github.com/drduh/gone/version"
)

// Load returns the configured application.
func Load() *App {
	app := App{}

	app.Debug = modeDebug
	app.Modes.Version = modeVersion

	s, err := settings.Load(pathConfig)
	if err != nil {
		log.Fatalf("failed loading settings: %v", err)
	}
	app.Settings = s

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

	auth.SetTarpit(app.TarpitDelay.Duration)

	app.Start()

	return &app
}
