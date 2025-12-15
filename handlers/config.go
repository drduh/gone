// Package handlers provides server operations
// on configured routes.
package handlers

const (
	formFieldClear     = "clear"
	formFieldDownloads = "downloads"
	formFieldDuration  = "duration"
	formFieldMessage   = "message"
	formFieldTheme     = "theme"
	formFieldWall      = "wall"
)

// Request contains server operation metadata.
type Request struct {

	// Path to handle ("/")
	Path string `json:"path"`

	// IP address including port ("127.0.0.1:12345")
	Address string `json:"address"`

	// Masked address ("User123")
	Mask string `json:"mask"`

	// User agent ("Mozilla/5.0 ...")
	Agent string `json:"agent"`
}
