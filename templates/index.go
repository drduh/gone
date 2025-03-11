package templates

import _ "embed"

//go:embed html/index.html
var HtmlIndex string

// Index HTML page elements
type Index struct {

	// Page title
	Title string

	// Upload handler path
	PathUpload string

	// Whether upload requires auth
	AuthUpload bool

	// Authentication header
	AuthHeader string

	// Authentication form field placeholder
	AuthHolder string
}
