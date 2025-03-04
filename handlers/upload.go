package handlers

import (
	"bytes"
	"io"
	"net/http"
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
			app.Log.Error("upload failed",
				"ip", ip, "ua", ua, "error", err.Error())
			return
		}
		defer file.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			writeJSON(w, http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
			app.Log.Error("upload failed",
				"ip", ip, "ua", ua, "error", err.Error())
			return
		}

		record := &config.File{
			Name:     handler.Filename,
			Uploaded: time.Now(),
			Size:     len(buf.Bytes()),
			Data:     buf.Bytes(),
			Owner: config.Owner{
				Address: ip,
				Agent:   ua,
			},
		}
		app.Storage.Files[record.Name] = record

		response := map[string]interface{}{
			"status": "ok",
			"data": map[string]interface{}{
				"name": record.Name,
				"size": record.Size,
				"time": record.Uploaded,
			},
		}

		writeJSON(w, http.StatusOK, response)
		app.Log.Info("upload complete",
			"ip", ip, "ua", ua,
			"name", record.Name, "size", record.Size)
	}
}
