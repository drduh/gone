package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	_ "embed"
)

//go:embed defaultSettings.json
var defaultSettings []byte

// Load returns the application configuration from
// a settings file or the default embedded settings.
func Load(path string) Settings {
	var settings Settings

	if err := loadSettings(defaultSettings, &settings); err != nil {
		log.Fatalf("failed loading default settings: %v", err)
	}

	if path != "" {
		if err := loadSettingsFromFile(path, &settings); err != nil {
			log.Fatalf("failed loading settings from file: %v", err)
		}
	}

	return settings
}

// loadSettings decodes a JSON byte slice into a Settings struct,
// rejecting unknown fields, trailing data, and invalid configurations.
func loadSettings(data []byte, settings *Settings) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(settings); err != nil {
		return fmt.Errorf("failed to decode settings: %w", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("failed to decode settings: trailing data")
		}
		return fmt.Errorf("failed to decode settings: trailing data: %w", err)
	}

	if err := settings.Validate(); err != nil {
		return fmt.Errorf("failed to validate settings: %w", err)
	}

	return nil
}

// loadSettingsFromFile loads settings from a file at path.
func loadSettingsFromFile(path string, settings *Settings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %w", path, err)
	}
	return loadSettings(data, settings)
}
