package templates

import (
	"embed"

	"github.com/drduh/gone/config"
)

//go:embed data/*.tmpl
var All embed.FS

// Index HTML page elements
type Index struct {

	// Selected CSS theme
	Theme string

	// Server name
	Hostname string

	// Time since server start
	Uptime string

	// Application version/build information
	Version map[string]string

	// Form field placeholder for duration
	DefaultDuration string

	// Authentication configuration
	config.Auth

	// Page properties
	config.Index

	// Page restrictions and limits
	config.Limits

	// Configured route paths
	config.Paths

	// Configured storage (Files and Messages)
	config.Storage
}
