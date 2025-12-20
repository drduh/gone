package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/drduh/gone/settings"
)

// TestParseFormInt tests integer form values are parsed.
func TestParseFormInt(t *testing.T) {
	def := 1
	maximum := 100
	tests := []struct {
		name    string
		query   string
		field   string
		def     int
		want    int
		maximum int
	}{
		{"valid", "/?downloads=5", "downloads",
			def, 5, maximum},
		{"space", "/?downloads= 5  ", "downloads",
			def, 5, maximum},
		{"missing", "/", "downloads",
			def, 1, maximum},
		{"invalid", "/?downloads=none", "downloads",
			def, 1, maximum},
		{"zero", "/?downloads=0", "downloads",
			def, def, maximum},
		{"negative", "/?downloads=-1", "downloads",
			def, def, maximum},
		{"fraction", "/?downloads=3.5", "downloads",
			def, def, maximum},
		{"large", "/?downloads=101", "downloads",
			def, maximum, maximum},
		{"xlarge", "/?downloads=999999999999999999999",
			"downloads", def, def, maximum}, // overflows int64
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := parseFormInt(req, tc.field, tc.def, tc.maximum)
			if got != tc.want {
				t.Fatalf("parseFormInt(%q) = %d; want %d",
					tc.query, got, tc.want)
			}
		})
	}
}

// TestParseFormDuration tests duration form values are parsed.
func TestParseFormDuration(t *testing.T) {
	def := 1 * time.Hour
	maximum := 8 * 24 * time.Hour
	tests := []struct {
		name    string
		query   string
		field   string
		def     time.Duration
		want    time.Duration
		maximum time.Duration
	}{
		{"valid", "/?duration=1h30m", "duration",
			def, 90 * time.Minute, maximum},
		{"space", "/?duration= 15m ", "duration",
			def, 15 * time.Minute, maximum},
		{"missing", "/", "duration",
			def, def, maximum},
		{"invalid", "/?duration=none", "duration",
			def, def, maximum},
		{"zero", "/?duration=0s", "duration",
			def, def, maximum},
		{"negative", "/?duration=-1h", "duration",
			def, def, maximum},
		{"fraction", "/?duration=1.5h", "duration",
			def, 90 * time.Minute, maximum},
		{"large", "/?duration=9999h", "duration",
			def, maximum, maximum},
		{"xlarge", "/?duration=99999999999h", "duration",
			def, def, maximum}, // overflows int64
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := parseFormDuration(req, tc.field, tc.def,
				settings.Duration{Duration: tc.maximum})
			if got != tc.want {
				t.Fatalf("parseFormDuration(%q) = %v; want %v",
					tc.query, got, tc.want)
			}
		})
	}
}
