package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/drduh/gone/config"
)

// Accepts content uploads
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

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

		downloadLimit := 0
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

		response := map[string]interface{}{
			"status": "ok",
			"data": map[string]interface{}{
				"fileName":       record.Name,
				"fileSize":       record.Size,
				"uploadTime":     record.Uploaded,
				"limitDownloads": record.LimitDownloads,
			},
		}

		writeJSON(w, http.StatusOK, response)
		app.Log.Info("upload complete",
			"name", record.Name, "size", record.Size,
			"ip", ip, "ua", ua)
	}
}
