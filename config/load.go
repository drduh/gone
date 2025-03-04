package config

import (
	"os"
	"time"

	"github.com/drduh/gone/version"
)

func Load() *App {
	hostname, _ := os.Hostname()

	return &App{
		Version:  version.Version,
		Hostname: hostname,
		Start:    time.Now(),
	}
}
