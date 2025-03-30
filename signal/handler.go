package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/drduh/gone/config"
)

// Logs and exits on terminal signals
func Setup(app *config.App) {
	sigChan := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
	signal.Notify(sigChan, signals...)
	go func() {
		s := <-sigChan
		app.Stop(s.String())
	}()
}
