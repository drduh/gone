package config

import (
	"os"
	"time"
)

// Start records the application start time.
func (a *App) Start() {
	a.StartTime = time.Now()
}

// Stop logs uptime and exits the application.
func (a *App) Stop(reason string) {
	a.Log.Info("stopping application",
		"reason", reason, "uptime", a.Uptime())
	os.Exit(0)
}

// Uptime returns the rounded duration since app start.
func (a *App) Uptime() string {
	return time.Since(a.StartTime).Round(
		time.Second).String()
}
