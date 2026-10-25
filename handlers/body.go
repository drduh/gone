package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

const (
	// Bytes per form character: utf8.UTFMax, each %XX-encoded.
	formBytesPerChar = 12

	// Headroom for field names, token, multipart boundaries.
	formBytesExtra = 64 << 10

	// Floor transfer rate used to scale deadlines with size.
	minTransferRate = 1 << 20 // bytes per second
)

// maxFormBytes returns the body limit for non-upload forms,
// derived from the largest configured text limit.
func maxFormBytes(app *config.App) int64 {
	chars := max(app.WallLimits.LengthChars,
		app.MessageLimits.LengthChars)

	return int64(chars)*formBytesPerChar + formBytesExtra
}

// bodyLimit returns the byte limit for a request body.
func bodyLimit(app *config.App, r *http.Request) int64 {
	if r.Pattern == app.Upload {
		return app.GetMaxFileBytes()
	}

	return maxFormBytes(app)
}

// transferDeadline returns base from now, plus time to
// move n bytes at the minimum transfer rate.
func transferDeadline(base time.Duration, n int64) time.Time {
	extra := time.Duration(n/minTransferRate) * time.Second

	return time.Now().Add(base + extra)
}

// parseBody bounds and parses the request body before auth.
func parseBody(
	w http.ResponseWriter,
	r *http.Request,
	app *config.App,
	req *Request) bool {
	limit := bodyLimit(app, r)

	reason := app.WallLimit
	if r.Pattern == app.Upload {
		reason = app.FileSize
	}

	if r.ContentLength > limit {
		Deny(w, http.StatusRequestEntityTooLarge, reason)
		app.Log.Error(reason,
			"contentLength", r.ContentLength,
			"limit", limit,
			"user", req)
		return false
	}

	if app.Read.Duration > 0 {
		n := r.ContentLength
		if n < 0 {
			n = limit
		}
		_ = http.NewResponseController(w).SetReadDeadline(
			transferDeadline(app.Read.Duration, n))
	}

	r.Body = http.MaxBytesReader(w, r.Body, limit)

	err := r.ParseForm()
	if err == nil {
		//gosec:disable G120: memory bounded by MaxBytesReader
		err = r.ParseMultipartForm(limit)
	}
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		Deny(w, http.StatusRequestEntityTooLarge, reason)
		app.Log.Error(reason,
			"limit", limit,
			"user", req)
		return false
	}

	return true
}
