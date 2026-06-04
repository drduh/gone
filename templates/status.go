package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Status represents server state and configuration.
type Status struct {

	// Defaults configuration
	settings.Default `json:"defaultOptions"`

	// Limits configuration
	settings.Limit `json:"limits"`

	// Index page configuration
	settings.Index `json:"indexPage"`

	// Storage content total sizes
	storage.Sizes `json:"storageSizes"`

	// Formatted time since start ("3m45s")
	Uptime string `json:"uptime,omitempty"`

	// Server hostname ("system")
	Hostname string `json:"hostname,omitempty"`

	// IP address the server is listening on ("127.0.0.1")
	ServerAddr string `json:"addr,omitempty"`

	// TCP port the server is listening on (8080)
	ServerPort int `json:"port,omitempty"`

	// Application version and build information
	Version map[string]string `json:"buildInfo,omitempty"`
}
