package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Index contains the main HTML page elements.
type Index struct {

	// Whether to display build, uptime and version in footer
	ShowBuild bool

	// Selected CSS theme
	Theme string

	// Server name
	Hostname string

	// Time since server start
	Uptime string

	// Application version/build information
	Version map[string]string

	// Placeholder text for duration form field
	DefaultDuration string

	// Authentication configuration
	settings.Auth

	// Page properties
	settings.Index

	// Page restrictions and limits
	settings.Limits

	// Configured route paths
	settings.Paths

	// Configured storage (Files and Messages)
	storage.Storage
}
