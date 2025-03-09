package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

// Time-based file expiration
type Duration struct {

	// Wrap time.Duration for JSON parsing
	time.Duration
}

// Parse strings like "10", "10s", "30m", "24h", etc.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed parsing duration: %w", err)
	}

	if isNumeric(s) {
		s += "s"
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("failed parsing duration: %w", err)
	}
	d.Duration = duration

	return nil
}

// Returns true if string contains numbers only
func isNumeric(str string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(str)
}
