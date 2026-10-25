package handlers

import (
	"net/http"
	"strconv"
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

		response := app.ListMessages()

		formContent := r.FormValue("download")
		if formContent == "allMessages" {
			app.Log.Debug("downloading all messages",
				"count", len(response),
				"user", req)
			app.ServeMessages(w)

			return
		}

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		app.Log.Info("serving messages",
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

		msgLength := charCount(formContent)
		if msgLength > app.MessageLimits.LengthChars {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgLength))
			app.Log.Error(app.MsgLength,
				"limit", app.MessageLimits.LengthChars,
				"length", msgLength,
				"user", req)

			return
		}

		t := time.Now()
		message := storage.Message{
			Data:          formContent,
			Agent:         req.Agent,
			Mask:          req.AddressMask,
			UploadTime:    t,
			UploadTimeFmt: t.Format(app.TimeFormat),
		}

		if !app.AddMessage(&message, app.MessageLimits.MaxCount) {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgCount))
			app.Log.Error(app.MsgCount,
				"limit", app.MessageLimits.MaxCount,
				"user", req)

			return
		}

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

		app.ClearMessages()
		app.Log.Info("cleared messages",
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, "messages cleared")
	}
}

// parseCount reads the Message count path value.
func parseCount(r *http.Request) (int, bool) {
	count, err := strconv.Atoi(r.PathValue("count"))
	if err != nil || count < 1 {
		return 0, false
	}

	return count, true
}

// MessageGet handles requests to read a Message by Count.
func MessageGet(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		count, ok := parseCount(r)
		if !ok {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgNotFound))
			return
		}

		msg, found := app.GetMessage(count)
		if !found {
			writeJSON(w, http.StatusNotFound,
				errorJSON(app.MsgNotFound))
			return
		}

		app.Log.Info("serving message",
			"count", count,
			"user", req)
		writeJSON(w, http.StatusOK, msg)
	}
}

// MessageDelete handles requests to delete a Message by Count.
func MessageDelete(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := AuthRequest(w, r, app)
		if req == nil {
			return
		}

		count, ok := parseCount(r)
		if !ok {
			writeJSON(w, http.StatusBadRequest,
				errorJSON(app.MsgNotFound))
			return
		}

		if !app.DeleteMessage(count) {
			writeJSON(w, http.StatusNotFound,
				errorJSON(app.MsgNotFound))
			return
		}

		app.Log.Info("deleted message",
			"count", count,
			"user", req)

		if req.IsBrowser {
			toPath(w, r, app.Root)
			return
		}

		writeJSON(w, http.StatusOK, "message deleted")
	}
}
