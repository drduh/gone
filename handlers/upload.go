package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Accepts content uploads
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if app.Settings.Auth.Require.Upload && !auth.Basic(app.Settings.Auth.Basic, r) {
			writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
			app.Log.Error(errorDeny,
				"action", "upload",
				"ip", ip, "ua", ua)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, http.StatusBadRequest,
				map[string]string{"error": err.Error()})
			app.Log.Error("upload form failed",
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}
		defer file.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
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
		}

		record := &config.File{
			Name:           handler.Filename,
			Uploaded:       time.Now(),
			LimitDownloads: downloadLimit,
			Size:           len(buf.Bytes()),
			Data:           buf.Bytes(),
			Owner: config.Owner{
				Address: ip,
				Agent:   ua,
			},
		}
		app.Storage.Files[record.Name] = record

		response := templates.File{
			Name: record.Name,
			Size: record.Size,
			Owner: templates.Owner{
				Address:  record.Owner.Address,
				Agent:    record.Owner.Agent,
				Uploaded: record.Uploaded,
			},
		}

		writeJSON(w, http.StatusOK, response)
		app.Log.Info("upload complete",
			"name", record.Name, "size", record.Size,
			"ip", ip, "ua", ua)
	}
}
