package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/settings"
	"github.com/drduh/gone/storage"
)

// Upload handles requests to upload Files to Storage.
func Upload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		maxFileBytes := app.GetMaxFileBytes()
		contentLength := r.ContentLength

		if contentLength > maxFileBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge,
				errorJSON(app.FileSize))
			app.Log.Error(app.FileSize,
				"fileSizeMB", settings.BytesToMegabytes(
					contentLength),
				"user", req)
			return
		}

		app.CountFiles()
		maxTotalBytes := app.GetMaxTotalFilesBytes()
		totalBytes := int64(app.SizeFiles) + contentLength

		if maxTotalBytes > 0 && totalBytes > maxTotalBytes {
			writeJSON(w, http.StatusRequestEntityTooLarge,
				errorJSON(app.FileSizeAll))
			app.Log.Error(app.FileSizeAll,
				"totalSize", app.SizeFiles,
				"user", req)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxFileBytes)
		//gosec:disable G120: memory bounded by MaxBytesReader
		if err := r.ParseMultipartForm(maxFileBytes); err != nil {
			writeJSON(w, http.StatusInternalServerError,
				errorJSON(app.Copy))
			app.Log.Error("upload failed",
				"error", err.Error(),
				"user", req)
			return
		}

		downloadLimit := parseFormInt(r, formFieldDownloads,
			app.Downloads, app.FileLimits.MaxDownloads)
		app.Log.Debug("got downloads form value",
			formFieldDownloads, downloadLimit,
			"user", req)

		durationLimit := parseFormDuration(r, formFieldDuration,
			app.Expiration.Duration, app.FileLimits.MaxDuration)
		app.Log.Debug("got duration form value",
			formFieldDuration, durationLimit,
			"user", req)

		formFileContent := r.MultipartForm.File["file"]
		if formFileContent == nil {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.Form))
			app.Log.Error(app.Form,
				"user", req)
			return
		}

		fileCount := len(formFileContent)
		if fileCount > app.FileLimits.MaxSelectFiles {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.UploadCount))
			app.Log.Error(app.UploadCount,
				"limit", app.FileLimits.MaxSelectFiles,
				"received", fileCount,
				"user", req)
			return
		}

		uploads := make([]storage.File, 0, fileCount)
		for _, fileHeader := range formFileContent {
			f, err := readUpload(fileHeader)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError,
					errorJSON(app.Copy))
				app.Log.Error(app.Copy,
					"error", err.Error(),
					"user", req)
				return
			}

			f.Name = storage.SanitizeName(
				fileHeader.Filename,
				app.FileLimits.NameExtraChars,
				app.FileLimits.NameLength)

			t := time.Now()
			f.Address = req.Address
			f.Agent = req.Agent
			f.Mask = req.AddressMask
			f.Duration = durationLimit
			f.UploadTime = t
			f.UploadTimeFmt = t.Format(app.TimeFormat)
			f.Allow = downloadLimit
			f.Scan()

			app.Files[f.ID] = f

			uploads = append(uploads, storage.File{
				ID:                f.ID,
				Name:              f.Name,
				Size:              f.Size,
				Sum:               f.Sum,
				Type:              f.Type,
				Address:           f.Address,
				Mask:              f.Mask,
				Agent:             f.Agent,
				UploadTimeFmt:     f.UploadTimeFmt,
				DurationRemaining: f.Duration.String(),
				Allow:             f.Allow,
			})
		}

		app.Log.Info("files uploaded",
			"files", uploads,
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
		} else {
			writeJSON(w, http.StatusOK, uploads)
		}
	}
}

// readUpload reads one multipart file part into a File.
func readUpload(fh *multipart.FileHeader) (*storage.File, error) {
	part, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open part: %w", err)
	}
	defer func() { _ = part.Close() }()

	data, err := io.ReadAll(part)
	if err != nil {
		return nil, fmt.Errorf("read part: %w", err)
	}

	return &storage.File{Data: data}, nil
}
