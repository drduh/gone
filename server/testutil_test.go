package server

import (
	"log/slog"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// newTestApp sets up a configured App for tests,
// ignoring logging.
func newTestApp(files map[string]*storage.File) *config.App {
	app := config.Load()
	app.Log = slog.New(slog.DiscardHandler)
	app.Files = files
	return app
}
