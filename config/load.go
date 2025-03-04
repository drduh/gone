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

	// Flag override
	settings.Modes.Debug = modeDebug
	settings.Modes.Version = modeVersion

	auditor, err := audit.StartAuditor(settings.Modes.Debug)
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
