package config

import (
	"log/slog"
	"time"
)

// Global application configuration
type App struct {

	// Application version
	Version string

	// Server hostname
	Hostname string

	// Server start time
	Start time.Time

	// Structured logger/output
	Log *slog.Logger
}
