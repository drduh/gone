package config

import (
	"context"
	"os"
	"time"

	"github.com/drduh/gone/audit"
	"github.com/drduh/gone/settings"
)

var (
	funcExit = os.Exit
	funcNow  = time.Now
)

// fatal logs err as JSON with default settings and exits.
func fatal(msg string, err error) {
	defaults, _ := settings.Load("")

	auditor, auditErr := audit.Start(&audit.Config{
		TimeFormat: defaults.TimeFormat,
	})
	if auditErr == nil {
		auditor.Log.Error(msg, "error", err.Error())
	}

	funcExit(1)
}

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
