package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Status represents server state and configuration.
type Status struct {

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime,omitempty"`

	// Server hostname ("system")
	Hostname string `json:"hostname,omitempty"`

	// TCP port the server is listening on (8080)
	Port int `json:"port,omitempty"`

	// Application version and build information
	Version map[string]string `json:"buildInfo,omitempty"`

	// Defaults configuration
	settings.Default `json:"defaultOptions,omitempty"`

	// Limits configuration
	settings.Limit `json:"limits,omitempty"`

	// Index page configuration
	settings.Index `json:"indexPage,omitempty"`

	// Storage content total sizes
	storage.Sizes `json:"storageSizes,omitempty"`
}
