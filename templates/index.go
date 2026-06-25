package templates

import (
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Index represents index HTML page elements.
type Index struct {

	// Authentication configuration
	settings.Auth

	// Content sharing default limits
	settings.Default

	// Error messages
	settings.Error

	// Index page HTML properties
	settings.Index

	// Content limits
	settings.Limit

	// Server paths
	settings.Paths

	// Content storage
	storage.Storage

	// Server name
	Hostname string

	// Display build details in footer
	ShowBuild bool

	// Selected CSS theme
	Theme string

	// Time since server start
	Uptime string

	// Application version/build information
	Version map[string]string
}
