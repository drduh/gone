package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

// loadSettings decodes a JSON byte slice into Settings,
// rejecting unknown/trailing data and invalid configs.
func loadSettings(data []byte, settings *Settings) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(settings); err != nil {
		return fmt.Errorf(
			"failed to decode settings: %w", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("failed to decode: trailing data: %w", err)
		}
		return fmt.Errorf("failed to decode: trailing data: %w", err)
	}

	if err := settings.Validate(); err != nil {
		return fmt.Errorf("failed to validate settings: %w", err)
	}

	return nil
}

// loadSettingsFromFile loads settings from a file at path.
func loadSettingsFromFile(p string, s *Settings) error {
	dir := filepath.Dir(p)
	file := filepath.Base(p)

	f, err := os.OpenInRoot(dir, file)
	if err != nil {
		return fmt.Errorf("failed to read '%s' from '%s': %w",
			file, dir, err)
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read '%s' from '%s': %w",
			file, dir, err)
	}

	return loadSettings(data, s)
}
