// Package storage defines user uploaded content.
package storage

import (
	"net/http"
	"time"
)

// Storage contains file and text content uploaded by users.
type Storage struct {

	// Uploaded files
	Files map[string]*File

	// Submitted text messages
	Messages map[int]*Message

	// Shared edit content
	WallContent string
}

// File represents a user-uploaded file.
type File struct {

	// Provided filename
	Name string `json:"name,omitempty"`

	// File content
	Data []byte `json:"data,omitempty"`

	// Downloads information
	Downloads `json:"downloads,omitempty"`

	// File length (for Content-Length header)
	Length string `json:"length,omitempty"`

	// File size (human-readable string)
	Size string `json:"size,omitempty"`

	// File type (based on name extension)
	Type string `json:"type,omitempty"`

	// Uploader information
	Owner `json:"owner,omitempty"`

	// Timing information
	Time `json:"time,omitempty"`
}

// Message represents a user-submitted text message.
type Message struct {

	// Counter (to order messages)
	Count int

	// Message content
	Data string

	// Owner information
	Owner

	// Timing information
	Time
}

// Owner contains metadata about a user.
type Owner struct {

	// Remote IP address
	Address string `json:"address,omitempty"`

	// User Agent header
	Agent string `json:"agent,omitempty"`

	// Full HTTP headers
	Headers http.Header `json:"headers,omitempty"`
}

// Time represents user content time metadata.
type Time struct {

	// Duration of file lifetime
	Duration time.Duration `json:"duration,omitempty"`

	// Formatted duration of file lifetime
	Allow string `json:"allow,omitempty"`

	// Formatted duration until expiration
	Remain string `json:"remain,omitempty"`

	// Absolute upload datetime
	Upload time.Time `json:"upload,omitempty"`
}

// Downloads represents user content downloads metadata.
type Downloads struct {

	// Number of allowed downloads
	Allow int `json:"allow,omitempty"`

	// Remaining number of downloads to expiration
	Remain int `json:"remain,omitempty"`

	// Total number of downloads
	Total int `json:"total,omitempty"`
}
