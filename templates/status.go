package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Status contains the server status and configuration response.
type Status struct {

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime"`

	// Server hostname ("system")
	Hostname string `json:"hostname"`

	// TCP port the server is listening on (8080)
	Port int `json:"port"`

	// Application version and build information
	Version map[string]string `json:"version"`

	// Storage content total sizes
	storage.Sizes `json:"sizes"`

	// Defaults configuration
	settings.Default `json:"default"`

	// Limits configuration
	settings.Limit `json:"limit"`

	// Index configuration
	settings.Index `json:"index"`

	// Storage content owner information
	storage.Owner `json:"owner"`
}
