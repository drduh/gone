package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Message handles requests to read Messages.
func Message(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		app.CountMessages()

		formContent := r.FormValue("download")
		if formContent == "allMessages" {
			app.Log.Debug("downloading all messages",
				"count", app.NumMessages,
				"user", req)
			app.ServeMessages(w)

			return
		}

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		response := app.Messages
		app.Log.Info("serving message(s)",
			"count", len(response),
			"user", req)

		writeJSON(w, http.StatusOK, response)
	}
}

// MessageAdd adds a validated Message to Storage.
func MessageAdd(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		formContent := strings.TrimSpace(
			r.PostFormValue(formFieldMessage))
		if formContent == "" {
			return
		}

		app.CountMessages()

		msgLength := len(formContent)
		if msgLength > app.MessageLimits.LengthChars {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgLength))
			app.Log.Error(app.MsgLength,
				"limit", app.MessageLimits.LengthChars,
				"length", msgLength,
				"user", req)

			return
		}

		if app.NumMessages >= app.MessageLimits.MaxCount {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgCount))
			app.Log.Error(app.MsgCount,
				"limit", app.MessageLimits.MaxCount,
				"count", app.NumMessages,
				"user", req)

			return
		}

		t := time.Now()
		message := storage.Message{
			Count: app.NumMessages + 1,
			Data:  formContent,
			Owner: storage.Owner{
				Agent: req.Agent,
				Mask:  req.Mask,
			},
			Time: storage.Time{
				UploadTime:    t,
				UploadTimeFmt: t.Format(app.TimeFormat),
			},
		}

		app.Messages = append(app.Messages, &message)
		app.Log.Info("added message",
			"count", message.Count,
			"length", msgLength,
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, message)
	}
}

// MessageClear handles requests to clear Messages.
func MessageClear(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		app.Log.Debug("clearing messages",
			"count", app.NumMessages,
			"user", req)

		app.ClearMessages()
		app.Log.Info("cleared messages",
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, "message(s) cleared")
	}
}
