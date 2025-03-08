package handlers

import (
	"time"

	"github.com/drduh/gone/config"
)

// Returns false if request is throttled
func throttle(app *config.App) bool {
	if app.Settings.Limits.PerMinute <= 0 {
		return false
	}

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	app.Throttle.Lease.Lock()
	defer app.Throttle.Lease.Unlock()

	times := make([]time.Time, 0, len(app.Throttle.Times))
	for _, t := range app.Throttle.Times {
		if t.After(cutoff) {
			times = append(times, t)
		}
	}

	if len(times) >= app.Settings.Limits.PerMinute {
		return true
	}

	times = append(times, now)
	app.Throttle.Times = times

	return false
}
