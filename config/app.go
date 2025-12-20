package config

import (
	"log/slog"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// App represents the global application configuration.
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
	settings.Settings

	// Uploaded content storage
	storage.Storage

	// Global rate limiting requests throttle
	auth.RequestThrottle
}
