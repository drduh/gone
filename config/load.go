package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	_ "embed"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/version"
)

//go:embed defaultSettings.json
var defaultSettings []byte

// Returns loaded application configuration
func Load() *App {
	app := App{}
	app.Modes.Debug = modeDebug
	app.Modes.Version = modeVersion

	var settings Settings
	if err := json.Unmarshal(defaultSettings, &settings); err != nil {
		log.Fatalf("failed to load default settings: %v", err)
	}
	if pathConfig != "" {
		settings = loadFromFile(pathConfig)
	}
	app.Settings = settings

	auditor, err := audit.StartAuditor(&audit.Config{
		Debug:      app.Modes.Debug,
		TimeFormat: settings.Audit.TimeFormat,
		Filename:   settings.Audit.Filename,
	})
	if err != nil {
		log.Fatalf("failed to start auditor: %v", err)
	}
	app.Log = auditor.Log

	app.Hostname = getHostname()
	app.Version = version.Short()
	app.Start = time.Now()
	app.Storage = Storage{
		Files:    make(map[string]*File),
		Messages: make(map[int]*Message),
	}

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

// Returns OS hostname or exits with error
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}
	return hostname
}
