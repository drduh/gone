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
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}

	var settings Settings
	if err := json.Unmarshal(defaultSettings, &settings); err != nil {
		log.Fatalf("failed to load default settings: %v", err)
	}

	if pathConfig != "" {
		settings = loadFromFile(pathConfig)
	}
	settings.Modes.Debug = modeDebug
	settings.Modes.Version = modeVersion

	auditor, err := audit.StartAuditor(&audit.Config{
		Debug:      settings.Modes.Debug,
		TimeFormat: settings.Audit.TimeFormat,
		Filename:   settings.Audit.Filename,
	})
	if err != nil {
		log.Fatalf("failed to start auditor: %v", err)
	}

	return &App{
		Version:  version.Short(),
		Hostname: hostname,
		Start:    time.Now(),
		Log:      auditor.Log,
		Settings: settings,
		Storage:  Storage{Files: make(map[string]*File)},
	}
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
