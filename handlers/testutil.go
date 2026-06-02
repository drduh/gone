package handlers

import (
	"log/slog"

	"github.com/drduh/gone/config"
)

const (
	testAddrAndPort = "127.0.0.1:12345"
	formContentType = "application/x-www-form-urlencoded"
	testContentMsgs = "hello, world!"
	testContentWall = "hello,\r\nworld!\r\n"
)

// newTestApp sets up a config for testing,
// ignoring logging and rate limiting.
func newTestApp() *config.App {
	app := config.Load()
	app.Log = slog.New(slog.DiscardHandler)
	app.ReqsPerMinute = 1000
	return app
}
