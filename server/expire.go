package server

import (
	"time"

	"github.com/drduh/gone/config"
)

// expiryWorker starts the expiration loop
// using the configured ticker interval.
func expiryWorker(app *config.App) {
	ticker := time.NewTicker(
		app.FileLimits.ExpiryCheck.Duration)
	defer ticker.Stop()
	runExpiryLoop(app, ticker.C)
}

// runExpiryLoop expires Files on each received tick.
func runExpiryLoop(app *config.App, t <-chan time.Time) {
	for range t {
		expireFiles(app)
	}
}

// expireFiles removes expired Files from Storage.
func expireFiles(app *config.App) {
	for _, f := range app.Files {
		lifetime := f.GetLifetime()

		app.Log.Debug("checking expiration",
			"id", f.ID,
			"name", f.Name,
			"allowed", f.Duration.String(),
			"available", lifetime.String(),
			"remaining", f.TimeRemaining().String())

		reason := f.IsExpired()
		if reason != "" {
			app.Expire(f)
			app.Log.Info("removed expired file",
				"reason", reason,
				"id", f.ID,
				"name", f.Name,
				"available", lifetime.String(),
				"downloads", f.Count)
		}
	}
}
