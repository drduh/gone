package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/drduh/gone/util"
)

// Time-based file expiration
type Duration struct {

	// Wrap time.Duration for JSON parsing
	time.Duration
}

// UnmarshalJSON parses strings like "10", "10s", "30m", "24h".
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed parsing duration json: %w", err)
	}

	if util.IsNumeric(s) {
		s += "s"
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("failed parsing duration string: %w", err)
	}
	d.Duration = duration

	return nil
}

// GetDuration converts a Duration to time.Duration.
func (d *Duration) GetDuration() time.Duration {
	return time.Duration(d.Duration)
}
