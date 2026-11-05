package config

import (
	"context"
	"log"
	"os"
	"time"
)

var (
	funcExit  = os.Exit
	funcFatal = log.Fatalf
	funcNow   = time.Now
)

// Start records the application start time and
// initiates Storage.
func (a *App) Start() {
	a.StartTime = funcNow()
	a.ClearStorage()
}

// Stop records uptime, drains the server within the
// shutdown timeout, clears storage and exits.
func (a *App) Stop(reason string) {
	if a.Log != nil {
		a.Log.Info("stopping application",
			"reason", reason,
			"uptime", a.Uptime())
	}

	if a.StopServer != nil {
		ctx, cancel := context.WithTimeout(
			context.Background(), a.Shutdown.Duration)
		defer cancel()

		if err := a.StopServer(ctx); err != nil &&
			a.Log != nil {
			a.Log.Warn("server stop incomplete",
				"error", err.Error())
		}
	}

	a.ClearStorage()
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
