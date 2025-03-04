package audit

import (
	"log"
	"log/slog"
	"os"
)

// Monitor and save application events
type Auditor struct {

	// Structured logger for use in app
	Log *slog.Logger
	slog.Handler
	*log.Logger
}

// Returns initialized Auditor ready for logging
func StartAuditor() (*Auditor, error) {
	opts := slog.HandlerOptions{}
	dest := os.Stdout
	handler := &Auditor{
		Handler: slog.NewJSONHandler(dest, &opts),
		Logger:  log.New(dest, "", 0),
	}
	return &Auditor{Log: slog.New(handler)}, nil
}
