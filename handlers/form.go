package handlers

import (
	"net/http"
	"strconv"
	"time"
)

// parseFormInt reads an integer form value or returns the default.
func parseFormInt(r *http.Request, field string, def int) int {
	input := r.FormValue(field)
	if input != "" {
		if v, err := strconv.Atoi(input); err == nil {
			return v
		}
	}
	return def
}

// parseFormDuration reads a duration form value or returns the default.
func parseFormDuration(r *http.Request, field string, def time.Duration) time.Duration {
	input := r.FormValue(field)
	if input != "" {
		if d, err := time.ParseDuration(input); err == nil {
			return d
		}
	}
	return def
}
