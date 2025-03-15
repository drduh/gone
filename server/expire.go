package server

import (
	"time"

	"github.com/drduh/gone/config"
)

// Checks Storage for expired files on a timer
func expiryWorker(app *config.App) {
	ticker := time.NewTicker(app.Settings.Limits.Ticker.Duration)
	defer ticker.Stop()

	for range ticker.C {
		for _, file := range app.Storage.Files {
			lifetime := time.Since(file.Time.Upload).Round(time.Second)
			app.Log.Debug("checking expiration",
				"filename", file.Name,
				"allowed", file.Time.Duration.String(),
				"available", lifetime.String(),
				"remaining", file.TimeRemaining().String())

			reason := file.IsExpired(app.Settings)
			if reason != "" {
				app.Storage.Expire(file)
				app.Log.Info("removed file",
					"reason", reason,
					"filename", file.Name,
					"available", lifetime.String(),
					"downloads", file.Downloads.Total)
			}
		}
	}
}
