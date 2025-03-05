package audit

import (
	"log"
	"log/slog"
	"os"
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

	// Structured logger for use in app
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

	// Message (short)
	Message string `json:"message"`

	// Event data (full)
	Data map[string]interface{} `json:"data"`
}

var cfg = &Config{}

// Returns initialized Auditor ready for logging
func StartAuditor(c *Config) (*Auditor, error) {
	cfg = c

	opts := slog.HandlerOptions{}
	if cfg.Debug {
		opts.Level = slog.LevelDebug
	}

	dest := os.Stdout

	handler := &Auditor{
		Handler: slog.NewJSONHandler(dest, &opts),
		Logger:  log.New(dest, "", 0),
	}

	return &Auditor{Log: slog.New(handler)}, nil
}
