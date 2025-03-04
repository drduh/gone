package signal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drduh/gone/config"
)

// Logs and exits on terminal signals
func Setup(app *config.App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-sigChan
		app.Log.Error("handled signal",
			"signal", s.String(),
			"uptime", time.Since(app.Start).String())
		os.Exit(0)
	}()
}
