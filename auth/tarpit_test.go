package auth

import (
	"testing"
	"time"
)

const base = 100 * time.Millisecond

// TestTarpitDelay tests configured delay duration.
func TestTarpitDelay(t *testing.T) {
	if testing.Short() {
		t.Skip("short test: skip tarpit delay")
	}

	cases := []struct {
		name    string
		tarpit  time.Duration
		wantMin time.Duration
		wantMax time.Duration
	}{
		{"zero does not delay", 0, 0,
			5 * time.Millisecond},
		{"negative does not delay", -1, 0,
			5 * time.Millisecond},
		{"delay is at least base", base, base,
			base + 25*time.Millisecond},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetTarpit(tc.tarpit)
			start := time.Now()
			ApplyTarpit()
			elapsed := time.Since(start)

			if elapsed < tc.wantMin {
				t.Errorf("elapsed %s < min %s",
					elapsed, tc.wantMin)
			}
			if elapsed > tc.wantMax {
				t.Errorf("elapsed %s > max %s",
					elapsed, tc.wantMax)
			}
		})
	}
}

// TestTarpitJitter tests varying delays.
func TestTarpitJitter(t *testing.T) {
	if testing.Short() {
		t.Skip("short test: skip tarpit jitter")
	}

	SetTarpit(base)

	const samples = 10
	delays := make([]time.Duration, samples)
	for i := range delays {
		start := time.Now()
		ApplyTarpit()
		delays[i] = time.Since(start)
	}

	minimum, maximum := delays[0], delays[0]
	for _, d := range delays[1:] {
		if d < minimum {
			minimum = d
		}
		if d > maximum {
			maximum = d
		}
	}

	if maximum-minimum < 500*time.Microsecond {
		t.Errorf("no jitter with %d samples (%s-%s)",
			samples, minimum, maximum)
	}
}
