package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drduh/gone/settings"
)

// parseFormInt reads an integer form value or returns the default.
func parseFormInt(r *http.Request, field string, def int, maximum int) int {
	input := r.FormValue(field)
	if input != "" {
		input = strings.TrimSpace(input)
		if v, err := strconv.ParseUint(input, 10, 64); err == nil {
			if v == 0 {
				return def
			}
			if v > uint64(maximum) {
				return maximum
			}
			return int(v)
		}
	}
	return def
}

// parseFormDuration reads a duration form value or returns the default.
func parseFormDuration(r *http.Request, field string,
	def time.Duration, maximum settings.Duration) time.Duration {
	input := r.FormValue(field)
	if input != "" {
		input = strings.TrimSpace(input)
		if d, err := time.ParseDuration(input); err == nil {
			if d < time.Second {
				return def
			}
			if d > maximum.GetDuration() {
				return maximum.GetDuration()
			}
			return d
		}
	}
	return def
}
