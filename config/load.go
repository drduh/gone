package config

import (
	"encoding/json"
	"fmt"
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
	if err := loadSettings(defaultSettings, &settings); err != nil {
		log.Fatalf("failed loading default settings: %v", err)
	}
	if pathConfig != "" {
		if err := loadSettingsFromFile(pathConfig, &settings); err != nil {
			log.Fatalf("failed loading settings from file: %v", err)
		}
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

// loadSettings unmarshals settings from a JSON byte slice
func loadSettings(data []byte, settings *Settings) error {
	if err := json.Unmarshal(data, settings); err != nil {
		return fmt.Errorf("failed to unmarshal settings: %w", err)
	}
	return nil
}

// loadSettingsFromFile loads settings from a file
func loadSettingsFromFile(path string, settings *Settings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}
	return loadSettings(data, settings)
}
