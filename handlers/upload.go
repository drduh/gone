package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
)

// Accepts content uploads
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if app.Settings.Auth.Require.Upload &&
			!auth.Basic(app.Settings.Auth.Header, app.Settings.Auth.Token, r) {
			writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
			app.Log.Error(errorDeny,
				"action", "upload",
				"ip", ip, "ua", ua)
			return
		}

		if throttle(app) {
			writeJSON(w, http.StatusTooManyRequests, responseErrorRateLimit)
			app.Log.Error(errorRateLimit,
				"action", "upload",
				"ip", ip, "ua", ua)
			return
		}

		maxBytes := int64(app.Settings.Limits.MaxSizeMb) << 20
		if r.ContentLength > maxBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge, responseErrorFileTooLarge)
			app.Log.Error(errorFileTooLarge,
				"sizeMb", r.ContentLength/(1<<20),
				"ip", ip, "ua", ua)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := r.ParseMultipartForm(maxBytes); err != nil {
			app.Log.Error("upload failed",
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}

		content, handler, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, http.StatusBadRequest,
				map[string]string{"error": err.Error()})
			app.Log.Error("upload form failed",
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}
		defer content.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, content); err != nil {
			writeJSON(w, http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
			app.Log.Error("upload copy failed",
				"error", err.Error(),
				"ip", ip, "ua", ua)
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
			Size: len(buf.Bytes()),
			Data: buf.Bytes(),
			Owner: config.Owner{
				Address: ip,
				Agent:   ua,
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
				Address: file.Owner.Address,
				Agent:   file.Owner.Agent,
			},
			Time: config.Time{
				Upload: file.Upload,
				Allow:  file.Time.Duration.String(),
			},
			Downloads: config.Downloads{
				Allow: file.Downloads.Allow,
			},
		}
		writeJSON(w, http.StatusOK, response)

		app.Log.Info("file uploaded",
			"filename", file.Name,
			"size", file.Size,
			"ip", ip, "ua", ua)
	}
}
