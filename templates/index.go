package templates

import (
	"embed"

	"github.com/drduh/gone/config"
)

//go:embed data/*.tmpl
var All embed.FS

// Index HTML page elements
type Index struct {

	// Page title
	Title string

	// Whether to allow theme selection
	ThemePick bool

	// CSS theme
	Theme string

	// Application version and build information
	Version     string
	VersionFull map[string]string

	// Configured route paths
	config.Paths

	// Authentication configuration
	config.Auth

	// Whether to require auth on routes
	AuthDownload bool
	AuthList     bool
	AuthMessage  bool
	AuthUpload   bool

	// Duration form field placeholder
	DefaultDuration string

	// Configured storage (files and messages)
	config.Storage
}
