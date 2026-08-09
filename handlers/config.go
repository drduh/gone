// Package handlers provides server operations
// on configured routes.
package handlers

const (
	formFieldClear     = "clear"
	formFieldDownload  = "download"
	formFieldDownloads = "downloads"
	formFieldDuration  = "duration"
	formFieldMessage   = "message"
	formFieldTheme     = "theme"
	formFieldWall      = "wall"

	templatesData = "data/*.tmpl"
)

// Request represents user request metadata.
type Request struct {

	// Request path ("/")
	Path string `json:"path"`

	// IP address including port ("127.0.0.1:12345")
	Address string `json:"address"`

	// Address mask ("User123")
	AddressMask string `json:"addressMask"`

	// User agent ("Mozilla/5.0 ...")
	Agent string `json:"agent"`

	// Whether the request originated from a browser
	IsBrowser bool `json:"isBrowser"`
}
