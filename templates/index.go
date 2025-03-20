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

	// CSS theme
	Theme string

	// Application version and build information
	Version     string
	VersionFull map[string]string

	// Route paths
	PathDownload string
	PathList     string
	PathMessage  string
	PathUpload   string

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
