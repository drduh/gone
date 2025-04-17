package templates

import (
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Index contains the main HTML page elements.
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
	storage.Storage
}
