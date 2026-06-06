package config

import (
	"os"
	"time"
)

var (
	funcExit = os.Exit
	funcNow  = time.Now
)

// Start records the application start time.
func (a *App) Start() {
	a.StartTime = funcNow()
}

// Stop records uptime and exits the application.
func (a *App) Stop(reason string) {
	if a.Log != nil {
		a.Log.Info("stopping application",
			"reason", reason,
			"uptime", a.Uptime())
	}
	funcExit(0)
}

// Uptime returns the rounded duration since start.
func (a *App) Uptime() string {
	now := funcNow()
	if a.StartTime.IsZero() ||
		now.Before(a.StartTime) {
		return "0s"
	}
	return now.Sub(a.StartTime).Round(
		time.Second).String()
}
