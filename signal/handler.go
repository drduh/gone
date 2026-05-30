// Package signal handles OS signals received by the application
// (such as SIGTERM), gracefully handling and logging them.
package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/drduh/gone/config"
)

// Setup configures and listens for signals.
func Setup(app *config.App) {
	sigChan := make(chan os.Signal, 2)
	signals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
	}
	signal.Notify(sigChan, signals...)

	go func() {
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGUSR1:
				app.Log.Info("clearing storage",
					"signal", s)
				app.ClearStorage()
			default:
				app.Stop(s.String())
				return
			}
		}
	}()
}
