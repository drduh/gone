package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "embed"
)

//go:embed defaultSettings.json
var defaultSettings []byte

// Get returns the loaded application configuration
// from path or embedded default settings file.
func Get(pathConfig string) Settings {
	var settings Settings
	if err := loadSettings(defaultSettings, &settings); err != nil {
		log.Fatalf("failed loading default settings: %v", err)
	}

	if pathConfig != "" {
		if err := loadSettingsFromFile(pathConfig, &settings); err != nil {
			log.Fatalf("failed loading settings from file: %v", err)
		}
	}

	return settings
}

// loadSettings unmarshals settings from a JSON byte slice.
func loadSettings(data []byte, settings *Settings) error {
	if err := json.Unmarshal(data, settings); err != nil {
		return fmt.Errorf("failed to unmarshal settings: %w", err)
	}
	return nil
}

// loadSettingsFromFile loads settings from a file at path.
func loadSettingsFromFile(path string, settings *Settings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}
	return loadSettings(data, settings)
}
