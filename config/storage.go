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
	Name string

	// Time of upload
	Uploaded time.Time

	// Size of file (in bytes)
	Size int

	// Number of downloads
	Downloads int

	// Raw file content
	Data []byte

	// File owner/uploader
	Owner
}

// Information about the uploader
type Owner struct {

	// Remote IP address
	Address string

	// User Agent header
	Agent string
}
