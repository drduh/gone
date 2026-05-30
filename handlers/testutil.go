package handlers

import (
	"io"
	"log/slog"

	"github.com/drduh/gone/config"
)

const (
	formContentType = "application/x-www-form-urlencoded"
	testContentMsgs = "hello, world!"
	testContentWall = "hello,\r\nworld!\r\n"
)

// newTestApp sets up an app config for testing.
func newTestApp() *config.App {
	app := config.Load()
	app.Log = slog.New(
		slog.NewTextHandler(io.Discard, nil))
	app.ReqsPerMinute = 1000
	return app
}
