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

	// Application modes (debug, version, etc.)
	Modes

	// Loaded and validated application settings
	Settings

	// Uploaded content storage
	Storage

	// Rate limit throttle for requests
	Throttle
}

// Application operation modes
type Modes struct {

	// Whether to display verbose debug output
	Debug bool

	// Whether to display version/build information
	Version bool
}
