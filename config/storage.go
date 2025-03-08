package config

import (
	"time"
)

// Storage for uploaded content
type Storage struct {

	// Collection of files
	Files map[string]*File
}

// An uploaded file
type File struct {

	// Provided filename
	Name string `json:"name,omitempty"`

	// Time of upload
	Uploaded time.Time `json:"uploaded,omitempty"`

	// Size of file (in bytes)
	Size int `json:"size,omitempty"`

	// Number of downloads
	Downloads int `json:"downloads,omitempty"`

	// User limit on number of downloads
	LimitDownloads int `json:"limitDownloads,omitempty"`

	// Raw file content
	Data []byte `json:"data,omitempty"`

	// Information about the uploader
	Owner `json:"owner,omitempty"`
}

// File owner information
type Owner struct {

	// Remote IP address
	Address string `json:"address,omitempty"`

	// User Agent header
	Agent string `json:"agent,omitempty"`
}

// Returns reason if file is expired
func (f *File) IsExpired() string {
	if f.Downloads >= f.LimitDownloads {
		return "limit downloads"
	}
	return ""
}
