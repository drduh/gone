package templates

import (
	"embed"

	"github.com/drduh/gone/config"
)

//go:embed data/*.tmpl
var All embed.FS

// Index HTML page elements
type Index struct {

	// Whether to allow theme selection
	ThemePick bool

	// CSS theme
	Theme string

	// Application version and build information
	Version     string
	VersionFull map[string]string

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

	// Configured storage (files and messages)
	config.Storage
}
