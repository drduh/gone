package handlers

import (
	"time"

	"github.com/drduh/gone/config"
)

// Returns true if request is allowed; not throttled by rate limit
func throttle(app *config.App) bool {
	if app.Settings.Limits.UploadsPM <= 0 {
		return true
	}

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	app.Storage.Throttle.Lease.Lock()
	defer app.Storage.Throttle.Lease.Unlock()

	validTimes := make([]time.Time, 0, len(app.Storage.Throttle.Times))
	for _, t := range app.Storage.Throttle.Times {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= app.Settings.Limits.UploadsPM {
		return false
	}

	validTimes = append(validTimes, now)
	app.Storage.Throttle.Times = validTimes

	return true
}
