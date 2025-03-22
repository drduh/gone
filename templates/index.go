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
	Path config.Paths

	// Whether routes require auth
	AuthDownload bool
	AuthList     bool
	AuthMessage  bool
	AuthUpload   bool

	// Authentication header
	AuthHeader string

	// Authentication form field placeholder
	AuthHolder string

	// Duration form field placeholder
	DefaultDuration string

	// Uploaded files
	Files map[string]*config.File

	// Text messages
	Messages map[int]*config.Message
}
