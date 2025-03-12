package templates

import _ "embed"

//go:embed data/index.html
var HtmlIndex string

// Index HTML page elements
type Index struct {

	// Page title
	Title string

	// Route paths
	PathDownload  string
	PathList      string
	PathUpload    string
	PathHeartbeat string

	// Whether routes require auth
	AuthDownload bool
	AuthList     bool
	AuthUpload   bool

	// Authentication header
	AuthHeader string

	// Authentication form field placeholder
	AuthHolder string
}
