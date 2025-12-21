// Package storage provides content storage and functions.
package storage

import (
	"net/http"
	"time"
)

const (
	storageVersion = "1"

	filenameMessages = "messages.txt"
	filenameWall     = "wall.txt"
)

// Storage represents content uploaded by users.
type Storage struct {

	// Uploaded files
	Files map[string]*File `json:"files,omitempty"`

	// Posted text messages
	Messages map[int]*Message `json:"messages,omitempty"`

	// Shared edit content
	WallContent string `json:"wallContent,omitempty"`

	// Storage content total sizes
	Sizes `json:"storageSizes,omitempty"`
}

// File represents a user-uploaded file.
type File struct {

	// Unique identifier
	Id string `json:"id,omitempty"`

	// Provided filename
	Name string `json:"name,omitempty"`

	// Contents of upload
	Data []byte `json:"data,omitempty"`

	// Content hash sum
	Sum string `json:"sum,omitempty"`

	// Downloads information
	Downloads `json:"downloads,omitempty"`

	// Number of bytes
	Bytes int `json:"bytes,omitempty"`

	// File length (for Content-Length header)
	Length string `json:"length,omitempty"`

	// File size (formatted string)
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

	// Message count/order
	Count int `json:"count,omitempty"`

	// Message content
	Data string `json:"data,omitempty"`

	// Owner information
	Owner `json:"owner,omitempty"`

	// Timing information
	Time `json:"time,omitempty"`
}

// Owner represents metadata about a user.
type Owner struct {

	// IP address with port
	Address string `json:"address,omitempty"`

	// Masked IP address
	Mask string `json:"mask,omitempty"`

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

// Sizes represents Storage content sizes.
type Sizes struct {

	// Number of characters in all Messages
	CharsMessages int `json:"charsMessages,omitempty"`

	// Number of characters in Wall content
	CharsWall int `json:"charsWall,omitempty"`

	// Number of lines in Wall content
	LinesWall int `json:"linesWall,omitempty"`

	// Number of Files
	NumFiles int `json:"numFiles,omitempty"`

	// Number of Messages
	NumMessages int `json:"numMessages,omitempty"`

	// Total size of all Files
	SizeFiles int `json:"sizeFiles,omitempty"`

	// Formatted total size of all Files
	SizeFilesFmt string `json:"sizeFilesFmt,omitempty"`
}
