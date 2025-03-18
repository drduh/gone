package templates

import (
	_ "embed"

	"github.com/drduh/gone/config"
)

//go:embed data/index.html
var HtmlIndex string

// Index HTML page elements
type Index struct {

	// Page title
	Title string

	// Application identifier
	Version string

	// Route paths
	PathDownload  string
	PathList      string
	PathUpload    string
	PathHeartbeat string

	// Whether routes require auth
	AuthDownload bool
	AuthList     bool
	AuthUpload   bool

	// Whether to style HTML with CSS
	Style bool

	// Authentication header
	AuthHeader string

	// Authentication form field placeholder
	AuthHolder string

	// Duration form field placeholder
	DefaultDuration string

	// Text messages
	Messages map[int]*config.Message
}
