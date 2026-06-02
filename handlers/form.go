package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drduh/gone/settings"
)

// parseFormInt reads an integer form value or returns the default.
func parseFormInt(r *http.Request, field string, def, maximum int) int {
	input := strings.TrimSpace(r.FormValue(field))
	if input == "" {
		return def
	}
	v, err := strconv.Atoi(input)
	if err != nil {
		return def
	}
	if v <= 0 {
		return def
	}
	if v > maximum {
		return maximum
	}
	return v
}

// parseFormDuration reads a duration form value or returns the default.
func parseFormDuration(r *http.Request, field string,
	def time.Duration, maximum settings.Duration) time.Duration {
	input := strings.TrimSpace(r.FormValue(field))
	if input == "" {
		return def
	}
	d, err := time.ParseDuration(input)
	if err != nil {
		return def
	}
	if d < time.Second {
		return def
	}
	maxDuration := maximum.GetDuration()
	if d > maxDuration {
		return maxDuration
	}
	return d
}
