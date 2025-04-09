package config

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// Storage for uploaded content
type Storage struct {

	// Uploaded files
	Files map[string]*File

	// Submitted text messages
	Messages map[int]*Message
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

// ClearStorage removes all Files and Messages from Storage
func (s *Storage) ClearStorage() {
	s.ClearFiles()
	s.ClearMessages()
}

// ClearFiles removes all Files from Storage
func (s *Storage) ClearFiles() {
	s.Files = make(map[string]*File)
}

// ClearMessages removes all Messages from Storage
func (s *Storage) ClearMessages() {
	s.Messages = make(map[int]*Message)
}

// CountFiles returns the number of Files in Storage
func (s *Storage) CountFiles() int {
	return len(s.Files)
}

// CountMessages returns the number of Messages in Storage
func (s *Storage) CountMessages() int {
	return len(s.Messages)
}

// Expire removes a File from Storage
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.Name)
}

// IsExpires returns a reason if File is expired
func (f *File) IsExpired(s Settings) string {
	if f.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if time.Since(f.Upload) > f.Duration {
		return "limit duration"
	}
	return ""
}

// NumRemaining returns the number of downloads remaining until expiration
func (f *File) NumRemaining() int {
	return f.Downloads.Allow - f.Total
}

// GetType returns File content type based on extension
func (f *File) GetType() string {
	t := mime.TypeByExtension(filepath.Ext(f.Name))
	if t == "" {
		t = "application/octet-stream"
	}
	return t
}

// TimeRemaining returns the relative duration remaining until expiration
func (f *File) TimeRemaining() time.Duration {
	return time.Until(
		f.Time.Upload.Add(f.Time.Duration)).Round(time.Second)
}

// FindFile Returns File, if found by name
func (s *Storage) FindFile(name string) *File {
	var file *File
	for _, f := range s.Files {
		if f.Name == name {
			file = f
			break
		}
	}
	return file
}

// Serve writes File as HTTP response
func (f *File) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Type", f.GetType())
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(f.Data)
	f.Total++
}

// ServeMessages writes all Messages as HTTP response
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=messages.txt")
	for _, msg := range s.Messages {
		_, err := fmt.Fprintf(w, "%d. %s\n", msg.Count, msg.Data)
		if err != nil {
			return
		}
	}
}
