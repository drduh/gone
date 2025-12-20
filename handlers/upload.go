package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Upload handles requests to upload Files to Storage.
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		maxFileBytes := app.GetMaxFileBytes()
		contentLength := r.ContentLength
		if contentLength > maxFileBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge, errorJSON(app.FileSize))
			app.Log.Error(app.FileSize,
				"fileSizeMb", contentLength/(1<<20), "user", req)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxFileBytes)
		if err := r.ParseMultipartForm(maxFileBytes); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Copy))
			app.Log.Error("upload failed", "error", err.Error(), "user", req)
			return
		}

		downloadLimit := parseFormInt(r,
			formFieldDownloads, app.Downloads, app.MaxDownloads)
		app.Log.Debug("got form value", formFieldDownloads, downloadLimit)

		durationLimit := parseFormDuration(r,
			formFieldDuration, app.Expiration.Duration, app.MaxDuration)
		app.Log.Debug("got form value", formFieldDuration, durationLimit)

		var upload storage.File
		var uploads []storage.File
		var wg sync.WaitGroup

		formFileContent := r.MultipartForm.File["file"]
		if formFileContent == nil {
			writeJSON(w, http.StatusBadRequest, errorJSON(app.Form))
			app.Log.Error(app.Form, "user", req)
			return
		}
		wg.Add(len(formFileContent))

		for _, fileHeader := range formFileContent {
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

				filename := storage.SanitizeName(fileHeader.Filename,
					app.MaxSizeName, app.AllowedSpecialChars)
				f := &storage.File{
					Name: filename,
					Data: buf.Bytes(),
					Owner: storage.Owner{
						Address: req.Address,
						Mask:    req.Mask,
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

				f.Scan()

				app.Files[f.Id] = f

				upload = storage.File{
					Id:   f.Id,
					Name: f.Name,
					Size: f.Size,
					Sum:  f.Sum,
					Type: f.Type,
					Owner: storage.Owner{
						Address: f.Address,
						Mask:    f.Mask,
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

		toRoot(w, r, app.Root)
		writeJSON(w, http.StatusOK, uploads)
		app.Log.Info("file(s) uploaded",
			"files", uploads, "user", req)
	}
}
