package server

import (
	"time"

	"github.com/drduh/gone/config"
)

// Removes files from Storage after duration limit
func expiryWorker(app *config.App) {
	period := app.Settings.Limits.Expiration.Duration
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for range ticker.C {
		for _, record := range app.Storage.Files {
			lifetime := time.Since(record.Uploaded)
			timeleft := period - lifetime
			app.Log.Debug("checking file expiration",
				"filename", record.Name,
				"duration", period.String(),
				"lifetime", lifetime.String(),
				"timeleft", timeleft.String())

			reason := record.IsExpired(app.Settings)
			if reason != "" {
				delete(app.Storage.Files, record.Name)
				app.Log.Info("removed file",
					"filename", record.Name,
					"reason", reason,
					"downloads", record.Downloads,
					"lifetime", lifetime.String())
			}
		}
	}
}
