// Package audit provides a JSON logger to monitor
// application events.
package audit

import (
	"log"
	"log/slog"
)

// Auditor represents a configured application logger.
type Auditor struct {

	// Configured options
	Config

	// Structured logger with custom handler
	Log *slog.Logger
	slog.Handler
	*log.Logger
}

// Config represents an Auditor configuration.
type Config struct {

	// Set logging level to Debug
	Debug bool

	// Format for time field
	TimeFormat string

	// File to write events to (stdout if unset)
	Filename string
}

// Event represents a unique Auditor log.
type Event struct {

	// Time of event ("Saturday Dec 13 12:00")
	Time string `json:"time"`

	// Severity level ("INFO", "DEBUG", etc.)
	Level string `json:"level"`

	// Summary message ("starting server")
	Message string `json:"message"`

	// Additional information ("data": {"port": 8080})
	Data map[string]interface{} `json:"data"`
}
