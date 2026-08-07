package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/drduh/gone/settings"
)

// TestParseFormInt tests parsing downloads form values.
func TestParseFormInt(t *testing.T) {
	t.Parallel()

	def := 1
	maximum := 100

	cases := []struct {
		name    string
		query   string
		want    int
	}{
		{"valid", "/?downloads=5", 5},
		{"space", "/?downloads= 5  ", 5},
		{"unencoded", "/?downloads=+5", 5},
		{"encoded", "/?downloads=%2B5", 5},
		{"leading zeros", "/?downloads=005", 5},
		{"padded", "/?downloads=%C2%A05%C2%A0", 5},
		{"missing", "/", def},
		{"invalid", "/?downloads=none", def},
		{"empty", "/?downloads=", def},
		{"zero", "/?downloads=0", def},
		{"negative", "/?downloads=-1", def},
		{"fraction", "/?downloads=3.5", def},
		{"hex", "/?downloads=0x5", def},
		{"trailing", "/?downloads=5abc", def},
		{"leading", "/?downloads=abc5", def},
		{"non ascii", "/?downloads=%EF%BC%95", def},
		{"exceed max", "/?downloads=1000", maximum},
		{"underscore", "/?downloads=1_000", def},
		{"with space", "/?downloads=1 000", def},
		{"only spaces", "/?downloads=   ", def},
		{"empty duplicate", "/?downloads=&downloads=5", def},
		{"valid duplicate", "/?downloads=5&downloads=10", 5},
		{"multiple", "/?duration=5m&downloads=10", 10},
		{"upper case", "/?Downloads=5", def},
		{"xlarge", "/?downloads=999999999999999999999",
			def}, // overflows int64
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(),
				http.MethodPost, tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := parseFormInt(
				req, formFieldDownloads, def, maximum)
			if got != tc.want {
				t.Fatalf("parseFormInt(%q) = %d; want %d",
					tc.query, got, tc.want)
			}
		})
	}
}

// TestParseFormDuration tests parsing duration form values.
func TestParseFormDuration(t *testing.T) {
	t.Parallel()

	def := 1 * time.Hour
	maximum := 8 * 24 * time.Hour

	cases := []struct {
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
		{"no-unit", "/?duration=3333", "duration",
			def, def, maximum},
		{"bad-unit", "/duration=8h", "duration",
			def, def, maximum},
		{"large", "/?duration=9999h", "duration",
			def, maximum, maximum},
		{"xlarge", "/?duration=99999999999h", "duration",
			def, def, maximum}, // overflows int64
		{"encoded", "/?duration=%32%34%68", "duration",
			def, 24 * time.Hour, maximum}, // "24h" encoded
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(),
				http.MethodPost, tc.query, nil)
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
