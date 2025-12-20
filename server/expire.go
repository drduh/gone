package server

import (
	"time"

	"github.com/drduh/gone/config"
)

// expiryWorker checks for and expires Files
// on the configured "ticker" schedule.
func expiryWorker(app *config.App) {
	ticker := time.NewTicker(app.Ticker.Duration)
	defer ticker.Stop()
	for range ticker.C {
		expireFiles(app)
	}
}

// expireFiles removes expired Files from Storage.
func expireFiles(app *config.App) {
	for _, f := range app.Files {
		lifetime := f.GetLifetime()
		app.Log.Debug("checking expiration",
			"id", f.Id, "name", f.Name,
			"allowed", f.Duration.String(),
			"available", lifetime.String(),
			"remaining", f.TimeRemaining().String())
		reason := f.IsExpired()
		if reason != "" {
			app.Expire(f)
			app.Log.Info("removed file", "reason", reason,
				"id", f.Id, "name", f.Name,
				"available", lifetime.String(),
				"downloads", f.Total)
		}
	}
}
