package audit

import (
	"log"
	"log/slog"

	"github.com/drduh/gone/util"
)

// Returns initialized Auditor ready for logging
func Start(c *Config) (*Auditor, error) {
	opts := slog.HandlerOptions{}
	if c.Debug {
		opts.Level = slog.LevelDebug
	}

	dest, err := util.GetOutput(c.Filename)
	if err != nil {
		return nil, err
	}

	handler := &Auditor{
		Config:  *c,
		Handler: slog.NewJSONHandler(dest, &opts),
		Logger:  log.New(dest, "", 0),
	}

	return &Auditor{Log: slog.New(handler)}, nil
}
