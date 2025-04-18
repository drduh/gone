// Package storage defines user uploaded content.
package storage

import (
	"net/http"
	"time"
)

// Storage contains Files and Messages from users
type Storage struct {

	// Uploaded files
	Files map[string]*File

	// Submitted text messages
	Messages map[int]*Message

	// Shared edit content
	WallContent string
}

// An uploaded file
type File struct {

	// Provided filename
	Name string `json:"name,omitempty"`

	// File content
	Data []byte `json:"data,omitempty"`

	// Downloads information
	Downloads `json:"downloads,omitempty"`

	// File size (bytes parsed to string)
	Size string `json:"size,omitempty"`

	// Uploader information
	Owner `json:"owner,omitempty"`

	// Timing information
	Time `json:"time,omitempty"`
}

// A submitted text message
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

// File owner information
type Owner struct {

	// Remote IP address
	Address string `json:"address,omitempty"`

	// User Agent header
	Agent string `json:"agent,omitempty"`

	// Full HTTP headers
	Headers http.Header `json:"headers,omitempty"`
}

// Timing information
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

// Downloads information
type Downloads struct {

	// Number of allowed downloads
	Allow int `json:"allow,omitempty"`

	// Remaining number of downloads to expiration
	Remain int `json:"remain,omitempty"`

	// Total number of downloads
	Total int `json:"total,omitempty"`
}
