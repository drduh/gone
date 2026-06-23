package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "embed"
)

//go:embed defaultSettings.json
var defaultSettings []byte

var errTrailingData = errors.New("trailing data in settings")

// Load returns the application configuration from
// a settings file or the default embedded settings.
func Load(path string) (Settings, error) {
	var settings Settings

	if err := loadSettings(
		defaultSettings, &settings); err != nil {
		return Settings{}, fmt.Errorf("default settings: %w",
			err)
	}

	if path != "" {
		if err := loadSettingsFromFile(
			path, &settings); err != nil {
			return Settings{}, fmt.Errorf("settings file: %w",
				err)
		}
	}

	return settings, nil
}

// loadSettings decodes a JSON byte slice into Settings,
// rejecting unknown/trailing data and invalid configs.
func loadSettings(data []byte, settings *Settings) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(settings); err != nil {
		return fmt.Errorf("failed to decode settings: %w",
			err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errTrailingData
	}

	if err := settings.Validate(); err != nil {
		return fmt.Errorf("failed to validate: %w",
			err)
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
