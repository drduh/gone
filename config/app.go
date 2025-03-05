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

	// Loaded and validated application settings
	Settings

	// Uploaded content storage
	Storage
}

// Application settings
type Settings struct {

	// Application mode
	Modes `json:"modes"`

	// Paths to route
	Paths `json:"paths"`

	// Limit routes
	Limits `json:"limits"`

	// TCP port to listen on
	Port int `json:"port"`

	// Auditor options
	Audit `json:"audit"`
}

// Auditor logging preferences
type Audit struct {

	// Optional file destination for logs
	Filename string `json:"logFile"`

	// Format for datetime in logs
	TimeFormat string `json:"timeFormat"`
}

// Application operation modes
type Modes struct {

	// Whether to display verbose debug output
	Debug bool

	// Whether to display version/build information
	Version bool
}

// Paths to route
type Paths struct {

	// File download ("/download")
	Download string `json:"download"`

	// Heartbeat/health check ("/heartbeat")
	Heartbeat string `json:"heartbeat"`

	// File list ("/list")
	List string `json:"list"`

	// Embedded/static file ("/static")
	Static string `json:"static"`

	// File upload ("/upload")
	Upload string `json:"upload"`
}

// Download and upload limits
type Limits struct {

	// Number of allowed downloads
	Downloads int `json:"downloads"`
}
