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

	// Application modes (debug, version, etc.)
	Modes

	// Auditor options
	Audit `json:"audit"`

	// Authentication requirements
	Auth `json:"auth"`

	// Limit routes
	Limits `json:"limits"`

	// Paths to route
	Paths `json:"paths"`

	// TCP port to listen on
	Port int `json:"port"`
}

// Application operation modes
type Modes struct {

	// Whether to display verbose debug output
	Debug bool

	// Whether to display version/build information
	Version bool
}

// Auditor logging preferences
type Audit struct {

	// Optional file destination for logs
	Filename string `json:"logFile"`

	// Format for datetime in logs
	TimeFormat string `json:"timeFormat"`
}

// Authentication requirements
type Auth struct {

	// String-based token/pass
	Basic string `json:"basic"`

	// Route authentication requirements
	Require struct {

		// Whether to require authentication for download
		Download bool `json:"download"`

		// Whether to require authentication for list
		List bool `json:"list"`

		// Whether to require authentication for upload
		Upload bool `json:"upload"`
	} `json:"require"`
}

// Download and upload limits
type Limits struct {

	// Number of allowed downloads
	Downloads int `json:"downloads"`
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
