package throttle

import "time"

// getCutoff returns the Throttle window cutoff (1 minute).
func getCutoff(t time.Time) time.Time {
	return t.Add(-1 * time.Minute)
}

// Allow returns true if the Throttle limit
// is not exceeded within the cutoff period.
func (t *Throttle) Allow(limit int) bool {
	if limit <= 0 {
		return true
	}

	now := time.Now()
	cut := getCutoff(now)

	t.Lease.Lock()
	defer t.Lease.Unlock()

	times := make([]time.Time, 0, len(t.RequestTimes))
	for _, t := range t.RequestTimes {
		if t.After(cut) {
			times = append(times, t)
		}
	}

	if len(times) >= limit {
		return false
	}

	times = append(times, now)
	t.RequestTimes = times

	return true
}
