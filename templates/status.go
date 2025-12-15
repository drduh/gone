package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Status contains the server status and configuration response.
type Status struct {

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime,omitempty"`

	// Server hostname ("system")
	Hostname string `json:"hostname,omitempty"`

	// TCP port the server is listening on (8080)
	Port int `json:"port,omitempty"`

	// Application version and build information
	Version map[string]string `json:"buildInfo,omitempty"`

	// Storage content total sizes
	storage.Sizes `json:"storageSizes,omitempty"`

	// Defaults configuration
	settings.Default `json:"defaultOptions,omitempty"`

	// Limits configuration
	settings.Limit `json:"limits,omitempty"`

	// Index configuration
	settings.Index `json:"indexPage,omitempty"`

	// Storage content owner information
	storage.Owner `json:"request,omitempty"`
}
