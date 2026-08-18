// Package storage provides content storage and functions.
package storage

import (
	"net/http"
	"time"
)

const (
	storageVersion = "1"

	filenameMsgs = "messages.txt"
	filenameWall = "wall.txt"

	expiryNone          ExpiryReason = ""
	expiryDownloadLimit ExpiryReason = "limit downloads"
	expiryDurationLimit ExpiryReason = "limit duration"
	invalidUploadTime   ExpiryReason = "invalid upload time"
)

// ExpiryReason identifies the reason for File expired.
type ExpiryReason string

// Storage represents content uploaded by users.
type Storage struct {

	// Storage content total sizes
	Sizes `json:"storageSizes"`

	// Wall metadata
	WallMeta `json:"wallMeta"`

	// Wall editable text area
	WallContent string `json:"wallContent"`

	// Uploaded files
	Files map[string]*File `json:"files,omitempty"`

	// Text messages
	Messages []*Message `json:"messages,omitempty"`
}

// File represents an uploaded file.
type File struct {

	// Uploader information
	Owner `json:"owner"`

	// Downloads information
	Downloads `json:"downloads"`

	// Timing information
	Time `json:"time"`

	// Unique identifier
	ID string `json:"id,omitempty"`

	// Provided filename
	Name string `json:"name,omitempty"`

	// File content (not encoded)
	Data []byte `json:"-"`

	// Content hash sum
	Sum string `json:"sum,omitempty"`

	// Number of bytes
	Bytes int `json:"bytes,omitempty"`

	// File length (for Content-Length header)
	Length string `json:"length,omitempty"`

	// File size (formatted string)
	Size string `json:"size,omitempty"`

	// File content type
	Type string `json:"type,omitempty"`

	// File content type display override
	TypeFmt string `json:"typeFmt,omitempty"`
}

// Message represents a user-submitted text message.
type Message struct {

	// Owner information
	Owner `json:"owner"`

	// Timing information
	Time `json:"time"`

	// Message count/order
	Count int `json:"count,omitempty"`

	// Message content
	Data string `json:"data,omitempty"`
}

// MessageParts represents extracted Message attributes.
type MessageParts struct {

	// Text content
	Text string

	// URL content
	URL string

	// Whether part contains URL
	HasURL bool
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

// Time represents content time information.
type Time struct {

	// Duration expiration
	Duration time.Duration `json:"duration,omitempty"`

	// Remaining duration until expiration
	DurationRemaining string `json:"durationRemaining,omitempty"`

	// Upload date and time (precise)
	UploadTime time.Time `json:"-"`

	// Formatted upload time
	UploadTimeFmt string `json:"uploadTimeFmt,omitempty"`
}

// Downloads represents user content downloads metadata.
type Downloads struct {

	// Number of allowed downloads
	Allow int `json:"allow,omitempty"`

	// Number of download requests
	Count int `json:"count,omitempty"`

	// Number of downloads remaining until expiration
	Remain int `json:"remain,omitempty"`
}

// Sizes represents Storage content sizes.
type Sizes struct {

	// Number of characters in all Messages
	CharsMessages int `json:"charsMessages,omitempty"`

	// Number of Files
	NumFiles int `json:"numFiles,omitempty"`

	// Number of Messages
	NumMessages int `json:"numMessages,omitempty"`

	// Total size of all Files
	SizeFiles int `json:"sizeFiles,omitempty"`

	// Formatted total size of all Files
	SizeFilesFmt string `json:"sizeFilesFmt,omitempty"`
}

// WallMeta represents Wall metadata.
type WallMeta struct {

	// Number of characters in Wall content
	WallChars int `json:"wallChars,omitempty"`

	// Number of lines in Wall content
	WallLines int `json:"wallLines,omitempty"`

	// Time of last modification
	WallModifiedTime time.Time `json:"-"`

	// Formatted time of last modification
	WallModifiedTimeFmt string `json:"wallModifiedTimeFmt,omitempty"`

	// Last modification user
	WallModifiedUser string `json:"wallModifiedUser,omitempty"`
}
