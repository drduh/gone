package config

import (
	"log/slog"
	"time"

	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/throttle"
)

// Global application configuration
type App struct {

	// Application version and build info
	Version map[string]string

	// Server hostname
	Hostname string

	// Server start time
	StartTime time.Time

	// Structured logger/output
	Log *slog.Logger

	// Application modes (debug, version, etc.)
	Modes

	// Loaded and validated application settings
	Settings

	// Uploaded content storage
	storage.Storage

	// Rate limit throttle for requests
	throttle.Throttle
}
