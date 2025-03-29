package audit

import (
	"log"
	"log/slog"

	"github.com/drduh/gone/util"
)

var cfg = &Config{}

// Returns initialized Auditor ready for logging
func Start(c *Config) (*Auditor, error) {
	cfg = c

	opts := slog.HandlerOptions{}
	if cfg.Debug {
		opts.Level = slog.LevelDebug
	}

	dest, err := util.GetOutput(cfg.Filename)
	if err != nil {
		return nil, err
	}

	handler := &Auditor{
		Handler: slog.NewJSONHandler(dest, &opts),
		Logger:  log.New(dest, "", 0),
	}

	return &Auditor{Log: slog.New(handler)}, nil
}
