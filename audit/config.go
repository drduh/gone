package audit

import (
	"log"
	"log/slog"
)

// Monitor for application events
type Auditor struct {

	// Configured options
	Config

	// Structured logger with custom handler
	Log *slog.Logger
	slog.Handler
	*log.Logger
}

// Auditor configuration
type Config struct {

	// Whether to set debug logging level
	Debug bool

	// Format for time field
	TimeFormat string

	// File to write events to (stdout if unset)
	Filename string
}

// Audit Event
type Event struct {

	// Time of event
	Time string `json:"time"`

	// Severity level ("INFO", "DEBUG", etc.)
	Level string `json:"level"`

	// Short summary
	Message string `json:"message"`

	// Full event
	Data map[string]interface{} `json:"data"`
}
