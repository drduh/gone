package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Heartbeat contains the server status response.
type Heartbeat struct {

	// Application version and build information
	Version map[string]string `json:"version"`

	// Server hostname ("system")
	Hostname string `json:"hostname"`

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime"`

	// TCP port the server is listening on (8080)
	Port int `json:"port"`

	// Number of Files in Storage
	FileCount int `json:"files"`

	// Number of Messages in Storage
	MessageCount int `json:"messages"`

	// Length of Wall content in Storage
	WallCount int `json:"wall"`

	// Defaults configuration
	settings.Default `json:"default"`

	// Limits configuration
	settings.Limit `json:"limit"`

	// Index configuration
	settings.Index `json:"index"`

	// File owner information
	storage.Owner `json:"owner"`
}
