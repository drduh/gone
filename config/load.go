package config

import (
	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/util"
	"github.com/drduh/gone/version"
)

// Load returns the configured application.
func Load() *App {
	app := App{
		Debug:       modeDebug,
		ShowVersion: modeShowVersion,
	}

	s, err := settings.Load(pathConfig)
	if err != nil {
		funcFatal("failed loading settings: %v", err)
		return nil
	}
	app.Settings = s

	auditor, err := audit.Start(&audit.Config{
		Debug:      app.Debug || app.LogDebug,
		Filename:   app.LogFilename,
		TimeFormat: app.TimeFormat,
	})
	if err != nil {
		funcFatal("failed starting auditor: %v", err)
		return nil
	}
	app.Log = auditor.Log

	app.Hostname = util.GetHostname()
	app.Version = version.Get(version.Full)

	auth.SetTarpit(app.TarpitDelay.Duration)

	if authToken != "" {
		app.Basic.Token = authToken
	}
	if app.Basic.Token == "" {
		app.Basic.Token = util.GetRandom("mask")
	}

	app.Start()

	return &app
}
