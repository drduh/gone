package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Status contains the server status and configuration response.
type Status struct {

	// Application version and build information
	Version map[string]string `json:"version"`

	// Server hostname ("system")
	Hostname string `json:"hostname"`

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime"`

	// TCP port the server is listening on (8080)
	Port int `json:"port"`

	// Defaults configuration
	settings.Default `json:"default"`

	// Limits configuration
	settings.Limit `json:"limit"`

	// Index configuration
	settings.Index `json:"index"`

	// Storage content owner information
	storage.Owner `json:"owner"`

	// Storage content total sizes
	storage.Sizes `json:"sizes"`
}
