package templates

import "time"

// List of file records
type List struct {

	// Name of the file
	Name string `json:"name"`

	// Size of file
	Size int `json:"size"`

	// File owner/uploader
	Owner `json:"owner"`
}

// Information about the uploader
type Owner struct {

	// Upload time
	Uploaded time.Time `json:"uploaded"`

	// Remote IP address
	Address string `json:"address"`

	// User Agent header
	Agent string `json:"agent"`
}
