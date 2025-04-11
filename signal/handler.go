// Package signal handles OS signals received by the application
// (such as SIGTERM), gracefully handling and logging them.
package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/drduh/gone/config"
)

// Setup configures and listens for terminal signals,
// logging them before successfully exiting.
func Setup(app *config.App) {
	sigChan := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
	signal.Notify(sigChan, signals...)
	go func() {
		s := <-sigChan
		app.Stop(s.String())
	}()
}
