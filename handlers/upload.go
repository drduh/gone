package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/util"
)

// Upload handles requests to upload File(s) into Storage.
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if !app.Allow(app.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.RateLimit))
			app.Log.Error(app.RateLimit, "user", req)
			return
		}

		maxBytes := app.GetMaxBytes()
		if r.ContentLength > maxBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge, errorJSON(app.FileSize))
			app.Log.Error(app.FileSize,
				"sizeMb", r.ContentLength/(1<<20), "user", req)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := r.ParseMultipartForm(maxBytes); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Copy))
			app.Log.Error("upload failed", "error", err.Error(), "user", req)
			return
		}

		downloadLimit := app.Downloads
		downloadLimitInput := r.FormValue("downloads")
		if limit, err := strconv.Atoi(downloadLimitInput); err == nil {
			downloadLimit = limit
			app.Log.Debug("got form value", "downloads", downloadLimit)
		}

		durationLimit := app.Expiration.Duration
		durationLimitInput := r.FormValue("duration")
		if limit, err := time.ParseDuration(durationLimitInput); err == nil {
			durationLimit = limit
			app.Log.Debug("got form value", "duration", durationLimit.String())
		}

		var upload storage.File
		var uploads []storage.File
		var wg sync.WaitGroup

		files := r.MultipartForm.File["file"]
		if files == nil {
			writeJSON(w, http.StatusBadRequest, errorJSON(app.Form))
			app.Log.Error(app.Form, "user", req)
			return
		}
		wg.Add(len(files))

		for _, fileHeader := range files {
			go func(fileHeader *multipart.FileHeader) {
				defer wg.Done()
				file, err := fileHeader.Open()
				if err != nil {
					app.Log.Error(app.Copy, "error", err.Error(), "user", req)
				}
				defer func() {
					if err := file.Close(); err != nil {
						app.Log.Error(app.Form, "error", err.Error(), "user", req)
						return
					}
				}()

				var buf bytes.Buffer
				if _, err := io.Copy(&buf, file); err != nil {
					writeJSON(w, http.StatusInternalServerError, errorJSON(app.Copy))
					app.Log.Error(app.Copy, "error", err.Error(), "user", req)
					return
				}

				f := &storage.File{
					Name: fileHeader.Filename,
					Size: util.FormatSize(len(buf.Bytes())),
					Data: buf.Bytes(),
					Owner: storage.Owner{
						Address: req.Address,
						Agent:   req.Agent,
					},
					Time: storage.Time{
						Duration: durationLimit,
						Upload:   time.Now(),
					},
					Downloads: storage.Downloads{
						Allow: downloadLimit,
					},
				}
				app.Files[f.Name] = f

				upload = storage.File{
					Name: f.Name,
					Size: f.Size,
					Owner: storage.Owner{
						Address: f.Address,
						Agent:   f.Agent,
					},
					Time: storage.Time{
						Upload: f.Upload,
						Allow:  f.Duration.String(),
					},
					Downloads: storage.Downloads{
						Allow: f.Downloads.Allow,
					},
				}
				uploads = append(uploads, upload)
			}(fileHeader)
		}
		wg.Wait()

		http.Redirect(w, r, "/", http.StatusSeeOther)

		writeJSON(w, http.StatusOK, uploads)

		app.Log.Info("file(s) uploaded",
			"files", uploads, "user", req)
	}
}
