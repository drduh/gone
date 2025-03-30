package config

import (
	"log/slog"
	"os"
	"time"
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
	Storage

	// Rate limit throttle for requests
	Throttle
}

// Record start time
func (a *App) Start() {
	a.StartTime = time.Now()
}

// Exits application with uptime log
func (a *App) Stop(reason string) {
	a.Log.Info("stopping application",
		"reason", reason, "uptime", a.Uptime())
	os.Exit(0)
}

// Returns rounded duration since app start
func (a *App) Uptime() string {
	return time.Since(a.StartTime).Round(
		time.Second).String()
}
