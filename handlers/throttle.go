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

	app.Storage.Throttle.Lease.Lock()
	defer app.Storage.Throttle.Lease.Unlock()

	fileTimes := make([]time.Time, 0, len(app.Storage.Throttle.Times))
	for _, t := range app.Storage.Throttle.Times {
		if t.After(cutoff) {
			fileTimes = append(fileTimes, t)
		}
	}

	if len(fileTimes) >= app.Settings.Limits.PerMinute {
		return true
	}

	fileTimes = append(fileTimes, now)
	app.Storage.Throttle.Times = fileTimes

	return false
}
