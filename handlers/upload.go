package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/util"
)

// Accepts content uploads
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if !isAllowed(app, r) {
			deny(w, app, req)
			return
		}

		if !app.Allow(app.Settings.Limits.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.Error.RateLimit))
			app.Log.Error(app.Error.RateLimit, "user", req)
			return
		}

		maxBytes := int64(app.Settings.Limits.MaxSizeMb) << 20
		if r.ContentLength > maxBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge, errorJSON(app.Error.FileSize))
			app.Log.Error(app.Error.FileSize,
				"sizeMb", r.ContentLength/(1<<20), "user", req)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := r.ParseMultipartForm(maxBytes); err != nil {
			app.Log.Error("upload failed",
				"error", err.Error(), "user", req)
			return
		}

		content, handler, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorJSON(app.Error.Form))
			app.Log.Error("upload form failed",
				"error", err.Error(), "user", req)
			return
		}
		defer content.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, content); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Error.Copy))
			app.Log.Error("upload copy failed",
				"error", err.Error(), "user", req)
			return
		}

		downloadLimit := app.Settings.Limits.Downloads
		downloadLimitInput := r.FormValue("downloads")
		if limit, err := strconv.Atoi(downloadLimitInput); err == nil {
			downloadLimit = limit
			app.Log.Debug("got form value", "downloads", downloadLimit)
		}

		durationLimit := app.Settings.Limits.Expiration.Duration
		durationLimitInput := r.FormValue("duration")
		if limit, err := time.ParseDuration(durationLimitInput); err == nil {
			durationLimit = limit
			app.Log.Debug("got form value", "duration", durationLimit.String())
		}

		file := &config.File{
			Name: handler.Filename,
			Size: util.FormatSize(len(buf.Bytes())),
			Data: buf.Bytes(),
			Owner: config.Owner{
				Address: req.Address,
				Agent:   req.Agent,
			},
			Time: config.Time{
				Duration: durationLimit,
				Upload:   time.Now(),
			},
			Downloads: config.Downloads{
				Allow: downloadLimit,
			},
		}
		app.Storage.Files[file.Name] = file

		response := config.File{
			Name: file.Name,
			Size: file.Size,
			Owner: config.Owner{
				Address: file.Address,
				Agent:   file.Agent,
			},
			Time: config.Time{
				Upload: file.Upload,
				Allow:  file.Duration.String(),
			},
			Downloads: config.Downloads{
				Allow: file.Downloads.Allow,
			},
		}
		writeJSON(w, http.StatusOK, response)

		app.Log.Info("file uploaded",
			"filename", file.Name, "size", file.Size, "user", req)
	}
}
