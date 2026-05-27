package util

import (
	"testing"
	"time"
)

// TestIsDaytimeAt tests seasonal daytime detection.
func TestIsDaytimeAt(t *testing.T) {
	tests := []struct {
		name  string
		month time.Month
		hour  int
		want  bool
	}{
		{"spring early", time.April, 5, false},
		{"spring day", time.April, 12, true},
		{"spring late", time.April, 20, false},
		{"summer early", time.July, 4, false},
		{"summer day", time.July, 6, true},
		{"summer late", time.July, 22, false},
		{"fall early", time.October, 6, false},
		{"fall day", time.October, 13, true},
		{"fall late", time.October, 19, false},
		{"winter early", time.December, 7, false},
		{"winter day", time.December, 13, true},
		{"winter late", time.December, 17, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(
				2026, tt.month, 1, tt.hour, 0, 0, 0, time.UTC)
			got := IsDaytimeAt(testTime)
			if got != tt.want {
				t.Errorf("IsDaytimeAt(%v) = %v, want %v",
					testTime, got, tt.want)
			}
		})
	}
}
