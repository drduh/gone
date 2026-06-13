package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Message handles requests to read and modify Messages.
func Message(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}

		app.CountMessages()

		if r.Method == http.MethodPost {
			if r.FormValue(formFieldClear) != "" {
				app.Log.Debug("clearing messages",
					"count", app.NumMessages,
					"user", req)
				app.ClearMessages()
				app.Log.Info("cleared messages",
					"user", req)
			}

			if !addMessage(w, r, app, req) {
				return
			}

			if req.IsBrowser {
				toPath(w, r, app.Root)
				return
			}
		}

		if r.URL.Query().Get("download") == "all" {
			app.Log.Debug("downloading messages",
				"count", app.NumMessages,
				"user", req)
			app.ServeMessages(w)

			return
		}

		app.Log.Info("serving message(s)",
			"count", len(app.Messages),
			"user", req)
		writeJSON(w, http.StatusOK, app.Messages)
	}
}

// addMessage validates and appends a Message.
func addMessage(
	w http.ResponseWriter,
	r *http.Request,
	app *config.App,
	req *Request,
) bool {
	formContent := strings.TrimSpace(
		r.PostFormValue(formFieldMessage))
	if formContent == "" {
		return true
	}

	msgLength := len(formContent)
	if msgLength > app.MessageLimits.LengthChars {
		writeJSON(w, http.StatusBadRequest,
			errorJSON(app.MsgLength))
		app.Log.Error(app.MsgLength,
			"limit", app.MessageLimits.LengthChars,
			"length", msgLength,
			"user", req)

		return false
	}

	msgCount := len(app.Messages)
	if msgCount >= app.MessageLimits.MaxCount {
		writeJSON(w, http.StatusBadRequest,
			errorJSON(app.MsgCount))
		app.Log.Error(app.MsgCount,
			"limit", app.MessageLimits.MaxCount,
			"count", msgCount,
			"user", req)

		return false
	}

	message := storage.Message{
		Count: msgCount + 1,
		Data:  formContent,
		Owner: storage.Owner{
			Agent: req.Agent,
			Mask:  req.Mask,
		},
		Time: storage.Time{
			Allow: time.Now().Format(app.TimeFormat),
		},
	}

	app.Messages = append(app.Messages, &message)
	app.Log.Info("added message",
		"count", message.Count,
		"length", msgLength,
		"user", req)

	return true
}
