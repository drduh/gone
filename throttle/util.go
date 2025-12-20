package throttle

import "time"

// getCutoff returns the Throttle window cutoff (1 minute).
func getCutoff(t time.Time) time.Time {
	return t.Add(-1 * time.Minute)
}

// Allow returns true if the Requests limit
// is not exceeded within the cutoff period.
func (r *Requests) Authorize(limit int) bool {
	if limit <= 0 {
		return true
	}

	now := time.Now()
	cut := getCutoff(now)

	r.Lease.Lock()
	defer r.Lease.Unlock()

	times := make([]time.Time, 0, len(r.RequestTimes))
	for _, t := range r.RequestTimes {
		if t.After(cut) {
			times = append(times, t)
		}
	}

	if len(times) >= limit {
		return false
	}

	times = append(times, now)
	r.RequestTimes = times

	return true
}
