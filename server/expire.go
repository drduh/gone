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
	expired := app.ExpireFiles()

	app.Log.Debug("checked expiration",
		"removed", len(expired))

	for _, f := range expired {
		app.Log.Info("removed expired file",
			"id", f.ID,
			"name", f.Name,
			"reason", f.Reason,
			"allowed", f.Duration.String(),
			"available", f.Lifetime().String(),
			"downloads", f.Count)
	}
}
