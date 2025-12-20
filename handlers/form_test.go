package handlers

import (
	"net/http"
	"testing"
	"time"
)

// TestParseFormInt tests integer form values are parsed.
func TestParseFormInt(t *testing.T) {
	tests := []struct {
		name  string
		query string
		field string
		def   int
		want  int
	}{
		{"valid", "/?downloads=5", "downloads", 10, 5},
		{"missing", "/", "downloads", 10, 10},
		{"invalid", "/?downloads=none", "downloads", 10, 10},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := parseFormInt(req, tc.field, tc.def)
			if got != tc.want {
				t.Fatalf("parseFormInt(%q) = %d; want %d",
					tc.query, got, tc.want)
			}
		})
	}
}

// TestParseFormDuration tests duration form values are parsed.
func TestParseFormDuration(t *testing.T) {
	tests := []struct {
		name  string
		query string
		field string
		def   time.Duration
		want  time.Duration
	}{
		{"valid", "/?duration=1h30m", "duration",
			2 * time.Hour, 90 * time.Minute},
		{"missing", "/", "duration",
			2 * time.Hour, 2 * time.Hour},
		{"invalid", "/?duration=none", "duration",
			2 * time.Hour, 2 * time.Hour},
		{"seconds", "/?duration=123s", "duration",
			2 * time.Hour, 2*time.Minute + 3*time.Second},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := parseFormDuration(req, tc.field, tc.def)
			if got != tc.want {
				t.Fatalf("parseFormDuration(%q) = %v; want %v",
					tc.query, got, tc.want)
			}
		})
	}
}
