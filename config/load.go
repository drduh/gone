package config

import (
	"encoding/json"
	"log"
	"os"

	_ "embed"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/util"
	"github.com/drduh/gone/version"
)

//go:embed defaultSettings.json
var defaultSettings []byte

// Returns loaded application configuration
func Load() *App {
	app := App{}
	app.Start()
	app.Debug = modeDebug
	app.Modes.Version = modeVersion

	var settings Settings
	if err := json.Unmarshal(defaultSettings, &settings); err != nil {
		log.Fatalf("failed loading default settings: %v", err)
	}
	if pathConfig != "" {
		settings = loadFromFile(pathConfig)
	}
	app.Settings = settings

	auditor, err := audit.Start(&audit.Config{
		Debug:      app.Debug,
		TimeFormat: settings.TimeFormat,
		Filename:   settings.Filename,
	})
	if err != nil {
		log.Fatalf("failed starting auditor: %v", err)
	}
	app.Log = auditor.Log

	app.Hostname = util.GetHostname()
	app.Version = version.Full()
	app.Clear()

	return &app
}

// Returns application settings from file at path
func loadFromFile(path string) Settings {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %v", path, err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Fatalf("failed to load settings: %v", err)
	}

	return settings
}
