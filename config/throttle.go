package config

import (
	"sync"
	"time"
)

// Throttle requests by time
type Throttle struct {

	// Record times to rate limit
	Times []time.Time

	// Throttle lock
	Lease sync.Mutex
}

// Allow returns true if the Throttle limit is not exceeded.
func (t *Throttle) Allow(limit int) bool {
	if limit <= 0 {
		return true
	}

	now := time.Now()
	cut := getCutoff(now)

	t.Lease.Lock()
	defer t.Lease.Unlock()

	times := make([]time.Time, 0, len(t.Times))
	for _, t := range t.Times {
		if t.After(cut) {
			times = append(times, t)
		}
	}

	if len(times) >= limit {
		return false
	}

	times = append(times, now)
	t.Times = times

	return true
}

// getCutoff returns the Throttle window cutoff (1 minute).
func getCutoff(t time.Time) time.Time {
	return t.Add(-1 * time.Minute)
}
