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
		name  string
		query string
		want  int
	}{
		{"valid", "/?downloads=5", 5},
		{"space", "/?downloads= 5  ", 5},
		{"unencoded", "/?downloads=+5", 5},
		{"encoded", "/?downloads=%2B5", 5},
		{"leading-zeros", "/?downloads=005", 5},
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
		{"non-ascii", "/?downloads=%EF%BC%95", def},
		{"exceed max", "/?downloads=1000", maximum},
		{"underscore", "/?downloads=1_000", def},
		{"with-space", "/?downloads=1 000", def},
		{"only-spaces", "/?downloads=   ", def},
		{"empty-duplicate", "/?downloads=&downloads=5", def},
		{"valid-duplicate", "/?downloads=5&downloads=10", 5},
		{"upper-case", "/?Downloads=5", def},
		{"multiple", "/?duration=5m&downloads=10", 10},
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
		name  string
		query string
		want  time.Duration
	}{
		{"valid", "/?duration=1h30m", 90 * time.Minute},
		{"space", "/?duration= 15m ", 15 * time.Minute},
		{"missing", "/", def},
		{"invalid", "/?duration=none", def},
		{"zero", "/?duration=0s", def},
		{"negative", "/?duration=-1h", def},
		{"fraction", "/?duration=1.5h", 90 * time.Minute},
		{"no-unit", "/?duration=3333", def},
		{"bad-unit", "/duration=8h", def},
		{"large", "/?duration=9999h", maximum},
		{"xlarge", "/?duration=99999999999h",
			def}, // overflows int64
		{"encoded", "/?duration=%32%34%68",
			24 * time.Hour}, // "24h" encoded
		{"encoded-plus",
			"/?duration=1h%2B30m", def}, // "1h+30m"
		{"unicode-microseconds",
			"/?duration=1500%C2%B5s", def}, // below 1s min
		{"microseconds",
			"/?duration=1500us", def}, // below 1s min
		{"leading-decimal",
			"/?duration=.5s", def}, // below 1s min
		{"overflow-boundary",
			"/?duration=2562047h47m16.854775807s", maximum},
		{"overflow-one-nanosecond",
			"/?duration=2562047h47m16.854775808s", def},
		{"tabs-and-newlines",
			"/?duration=%09%0A15m%0D%0A", 15 * time.Minute},
		{"plus-as-space",
			"/?duration=+15m+", 15 * time.Minute},
		{"duplicate-first-valid",
			"/?duration=15m&duration=2h", 15 * time.Minute},
		{"duplicate-first-invalid",
			"/?duration=invalid&duration=15m", def},
		{"compound-fraction",
			"/?duration=1h0.5m", time.Hour + 30*time.Second},
		{"minimum", "/?duration=1s", time.Second},
		{"below-minimum", "/?duration=999ms", def},
		{"subsecond", "/?duration=1ns", def},
		{"at-maximum", "/?duration=192h", maximum},
		{"case-sensitive-unit", "/?duration=15M", def},
		{"just-over-maximum", "/?duration=192h1ns", maximum},
		{"unsupported-days", "/?duration=1d", def},
		{"plus-zero", "/?duration=+0s", def},
		{"negative-zero", "/?duration=-0s", def},
		{"mixed-sign", "/?duration=1h-30m", def},
		{"embedded-nul", "/?duration=15m%00ignored", def},
		{"html", "/?duration=%3Cscript%3Ealert(1)%3C%2Fscript%3E", def},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(),
				http.MethodPost, tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			got := parseFormDuration(req, formFieldDuration, def,
				settings.Duration{Duration: maximum})
			if got != tc.want {
				t.Fatalf("parseFormDuration(%q) = %v; want %v",
					tc.query, got, tc.want)
			}
		})
	}
}
