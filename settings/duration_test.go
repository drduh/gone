package settings

import (
	"encoding/json"
	"testing"
	"time"
)

// TestDurationUnmarshal tests interval string parsing.
func TestDurationUnmarshal(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
		err  bool
	}{
		{`"10"`, 10 * time.Second, false},
		{`"30m"`, 30 * time.Minute, false},
		{`"2h"`, 2 * time.Hour, false},
		{`"soon"`, 0, true},
		{`10`, 0, true},
	}

	for _, tc := range cases {
		var d Duration

		err := json.Unmarshal([]byte(tc.in), &d)
		if (err != nil) != tc.err {
			t.Fatalf("%s: err = %v; want error %v",
				tc.in, err, tc.err)
		}

		if d.GetDuration() != tc.want {
			t.Fatalf("%s: got %v; want %v",
				tc.in, d.Duration, tc.want)
		}
	}
}
