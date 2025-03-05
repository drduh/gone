package audit

import (
	"log"
	"log/slog"
)

// Auditor configuration
type Config struct {

	// Whether to output verbose debug messages
	Debug bool

	// Format for datetime in events
	TimeFormat string

	// Name of log file to write
	Filename string
}

// Monitor for application events
type Auditor struct {

	// Structured logger with custom handler
	Log *slog.Logger
	slog.Handler
	*log.Logger
}

// An audit event
type Event struct {

	// Time of event
	Time string `json:"time"`

	// Severity level ("INFO", "DEBUG", etc.)
	Level string `json:"level"`

	// Message summary (short)
	Message string `json:"message"`

	// Event data (full)
	Data map[string]interface{} `json:"data"`
}
