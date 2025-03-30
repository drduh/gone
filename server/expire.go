package server

import (
	"time"

	"github.com/drduh/gone/config"
)

// Runs expiration check on configured schedule
func expiryWorker(app *config.App) {
	ticker := time.NewTicker(app.Ticker.Duration)
	defer ticker.Stop()
	for range ticker.C {
		expireFiles(app)
	}
}

// Removes expired files from Storage
func expireFiles(app *config.App) {
	for _, f := range app.Files {
		lifetime := time.Since(f.Time.Upload).Round(time.Second)
		app.Log.Debug("checking expiration",
			"filename", f.Name,
			"allowed", f.Duration.String(),
			"available", lifetime.String(),
			"remaining", f.TimeRemaining().String())
		reason := f.IsExpired(app.Settings)
		if reason != "" {
			app.Expire(f)
			app.Log.Info("removed file",
				"reason", reason, "filename", f.Name,
				"available", lifetime.String(),
				"downloads", f.Total)
		}
	}
}
