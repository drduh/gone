package config

import (
	"net/http"
	"time"
)

// Storage for uploaded content
type Storage struct {

	// Uploaded files
	Files map[string]*File

	// Submitted text messages
	Messages map[int]*Message
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

// An uploaded file
type File struct {

	// Provided filename
	Name string `json:"name,omitempty"`

	// File size (bytes parsed into string)
	Size string `json:"size,omitempty"`

	// Uploader information
	Owner `json:"owner,omitempty"`

	// Timing information
	Time `json:"time,omitempty"`

	// Downloads information
	Downloads `json:"downloads,omitempty"`

	// File content
	Data []byte
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
func (f *File) TimeRemaining() time.Duration {
	return time.Until(
		f.Time.Upload.Add(f.Time.Duration)).Round(
		time.Second)
}

// Returns reason if File is expired
func (f *File) IsExpired(s Settings) string {
	if f.Downloads.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if time.Since(f.Time.Upload) > f.Time.Duration {
		return "limit duration"
	}
	return ""
}

// Removes File from Storage
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.Name)
}

// Clears Messages from Storage
func (s *Storage) ClearMessages() {
	s.Messages = make(map[int]*Message)
}

// Counts Messages in Storage
func (s *Storage) CountMessages() int {
	return len(s.Messages)
}
