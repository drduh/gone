// Package signal handles OS signals received by the application
// (such as SIGTERM), gracefully handling and logging them.
package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/drduh/gone/config"
)

const sigChanBuffer = 2

// Setup configures and listens for signals.
func Setup(app *config.App) {
	sigChan := make(chan os.Signal, sigChanBuffer)
	signals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	}
	signal.Notify(sigChan, signals...)

	go func() {
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGUSR1:
				app.CountStorage()
				app.Log.Info("clearing storage",
					"signal", s,
					"sizes", app.Sizes)
				app.ClearStorage()

			case syscall.SIGUSR2:
				app.Log.Info("auth token",
					"signal", s,
					"field", app.Basic.Field,
					"token", app.Basic.Token)

			default:
				app.Stop(s.String())
				return
			}
		}
	}()
}
