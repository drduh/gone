package auth

import "time"

const throttleInterval = 1 * time.Minute

// Authorize returns true if the RequestThrottle limit
// for addr is not exceeded within the cutoff period.
func (r *RequestThrottle) Authorize(addr string, limit int) bool {
	return r.isAuthorizedAt(addr, limit, time.Now())
}

// getCutoff returns the rate-limit window cutoff.
func getCutoff(t time.Time) time.Time {
	return t.Add(-throttleInterval)
}

// recent returns times after cut, in place.
func recent(times []time.Time, cut time.Time) []time.Time {
	kept := times[:0]
	for _, t := range times {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}

	return kept
}

// isAuthorizedAt returns true if the RequestThrottle limit
// for addr is not exceeded within the cutoff period.
func (r *RequestThrottle) isAuthorizedAt(
	addr string,
	limit int,
	t time.Time,
) bool {
	if limit <= 0 {
		return true
	}

	cut := getCutoff(t)

	r.Lease.Lock()
	defer r.Lease.Unlock()

	if r.Clients == nil {
		r.Clients = make(map[string][]time.Time)
	}

	if t.Sub(r.swept) >= throttleInterval {
		for a, times := range r.Clients {
			if times = recent(times, cut); len(times) == 0 {
				delete(r.Clients, a)
			} else {
				r.Clients[a] = times
			}
		}

		r.swept = t
	}

	times := recent(r.Clients[addr], cut)
	if len(times) >= limit {
		r.Clients[addr] = times
		return false
	}

	r.Clients[addr] = append(times, t)

	return true
}
