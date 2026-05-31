package audit

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/drduh/gone/util"
)

// Start returns an initialized Auditor for logging.
func Start(c *Config) (*Auditor, error) {
	dest, err := util.GetOutput(c.Filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to set up %s: %w", c.Filename, err)
	}

	opts := slog.HandlerOptions{}
	if c.Debug {
		opts.Level = slog.LevelDebug
	}

	handler := &Auditor{
		Config:  *c,
		Handler: slog.NewJSONHandler(dest, &opts),
		Logger:  log.New(dest, "", 0),
	}

	return &Auditor{Log: slog.New(handler)}, nil
}
