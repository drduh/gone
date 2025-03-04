package config

import (
	"os"
	"time"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/version"
)

func Load() *App {
	hostname, _ := os.Hostname()

	auditor, _ := audit.StartAuditor()

	return &App{
		Version:  version.Version,
		Hostname: hostname,
		Start:    time.Now(),
		Log:      auditor.Log,
	}
}
