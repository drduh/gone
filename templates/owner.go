package templates

import "time"

// File uploader/owner information
type Owner struct {

	// Upload time
	Uploaded time.Time `json:"uploaded"`

	// Remote IP address
	Address string `json:"address"`

	// User Agent header
	Agent string `json:"agent"`
}
