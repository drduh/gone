package audit

import (
	"bytes"
	"context"
	"log"
	"log/slog"
)

type testHandler struct{}

func (testHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h testHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h testHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h testHandler) WithGroup(string) slog.Handler {
	return h
}

// newTestAuditor creates a new Auditor writing to buf.
func newTestAuditor(buf *bytes.Buffer) *Auditor {
	a := &Auditor{
		Config: Config{
			TimeFormat: "2006-01-02 15:04:05",
		},
		Handler: testHandler{},
		Logger:  log.New(buf, "", 0),
	}
	a.Log = slog.New(a)
	return a
}
