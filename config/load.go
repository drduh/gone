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

	auditor, err := audit.StartAuditor()
	if err != nil {
		log.Fatalf("failed to start auditor: %v", err)
	}

	var s Settings
	if err := json.Unmarshal(defaultSettings, &s); err != nil {
		log.Fatalf("failed to load default settings: %v", err)
	}

	return &App{
		Version:  version.Short(),
		Hostname: hostname,
		Start:    time.Now(),
		Log:      auditor.Log,
		Settings: s,
	}
}
