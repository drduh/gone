package audit

import (
	"log"
	"log/slog"
	"os"
)

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
