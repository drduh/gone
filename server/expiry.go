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
			app.Log.Debug("checking file expiration",
				"duration", period.String(),
				"file", record.Name)

			reason := record.IsExpired(app.Settings)
			if reason != "" {
				delete(app.Storage.Files, record.Name)
				app.Log.Info("expired file",
					"reason", reason,
					"name", record.Name,
					"downloads", record.Downloads,
					"lifetime", time.Since(record.Uploaded).String())
			}
		}
	}
}
