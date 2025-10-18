package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Heartbeat contains the server status response.
type Heartbeat struct {

	// Application version
	Version map[string]string `json:"version"`

	// Server hostname
	Hostname string `json:"hostname"`

	// Time since start ("3m45s")
	Uptime string `json:"uptime"`

	// TCP port server is listening on
	Port int `json:"port"`

	// Number of Files in storage
	FileCount int `json:"files"`

	// Number of Messages in storage
	MessageCount int `json:"messages"`

	// Defaults configuration
	settings.Default `json:"default"`

	// Limits configuration
	settings.Limit `json:"limit"`

	// Index configuration
	settings.Index `json:"index"`

	// File owner information
	storage.Owner `json:"owner"`
}
