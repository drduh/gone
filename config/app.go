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

	// TCP port to listen on
	Port int `json:"port"`

	// Paths to route
	Paths `json:"paths"`
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
