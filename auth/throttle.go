package auth

import "time"

const throttleInterval = 1 * time.Minute

// Authorize returns true if the RequestThrottle limit
// is not exceeded within the cutoff period.
func (r *RequestThrottle) Authorize(limit int) bool {
	return r.isAuthorizedAt(limit, time.Now())
}

// getCutoff returns the rate-limit window cutoff.
func getCutoff(t time.Time) time.Time {
	return t.Add(-throttleInterval)
}

// isAuthorizedAt returns true if the RequestThrottle limit
// is not exceeded within the cutoff period.
func (r *RequestThrottle) isAuthorizedAt(
	limit int,
	t time.Time,
) bool {
	if limit <= 0 {
		return true
	}

	cut := getCutoff(t)

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

	times = append(times, t)
	r.RequestTimes = times

	return true
}
