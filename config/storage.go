package config

import (
	"net/http"
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

	// File size (in bytes)
	Size int `json:"size,omitempty"`

	// Raw file content
	Data []byte `json:"data,omitempty"`

	// Information about the uploader
	Owner `json:"owner,omitempty"`

	// Information about downloads
	Downloads `json:"downloads,omitempty"`
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

// File downloads information
type Downloads struct {

	// Number of allowed downloads
	Allow int `json:"allowed,omitempty"`

	// Remaining number of downloads to expiration
	Remain int `json:"remain,omitempty"`

	// Total number of downloads
	Total int `json:"total,omitempty"`
}

// Returns number of remaining allowed downloads
func (f *File) NumRemaining() int {
	return f.Downloads.Allow - f.Downloads.Total
}

// Returns reason if file is expired
func (f *File) IsExpired(s Settings) string {
	if f.Downloads.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if time.Since(f.Uploaded) > s.Limits.Expiration.Duration {
		return "limit duration"
	}
	return ""
}
