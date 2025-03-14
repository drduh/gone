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

	// File size (in bytes)
	Size int `json:"size,omitempty"`

	// Raw file content
	Data []byte `json:"data,omitempty"`

	// Uploader information
	Owner `json:"owner,omitempty"`

	// Timing information
	Time `json:"time,omitempty"`

	// Downloads information
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

// Timing information
type Time struct {

	// Formatted duration until expiration
	Remain string `json:"remain,omitempty"`

	// Absolute upload datetime
	Upload time.Time `json:"upload,omitempty"`
}

// Downloads information
type Downloads struct {

	// Number of allowed downloads
	Allow int `json:"allowed,omitempty"`

	// Remaining number of downloads to expiration
	Remain int `json:"remain,omitempty"`

	// Total number of downloads
	Total int `json:"total,omitempty"`
}

// Returns number of remaining downloads until expiration
func (f *File) NumRemaining() int {
	return f.Downloads.Allow - f.Downloads.Total
}

// Returns relative duration remaining until expiration
func (f *File) TimeRemaining(s Settings) time.Duration {
	return time.Until(
		f.Time.Upload.Add(s.Limits.Expiration.Duration)).Round(
		time.Second)
}

// Returns reason if file is expired
func (f *File) IsExpired(s Settings) string {
	if f.Downloads.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if time.Since(f.Time.Upload) > s.Limits.Expiration.Duration {
		return "limit duration"
	}
	return ""
}
