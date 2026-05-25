package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Message handles requests to read and modify Messages in Storage.
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

			message := storage.Message{
				Count: app.NumMessages,
				Owner: storage.Owner{
					Agent: req.Agent,
					Mask:  req.Mask,
				},
				Time: storage.Time{
					Allow: time.Now().Format(app.TimeFormat),
				},
			}

			formContent := strings.TrimSpace(
				r.PostFormValue(formFieldMessage))
			if formContent != "" {
				msgLength := len(formContent)
				if msgLength > app.MessageLimits.LengthChars {
					writeJSON(w, http.StatusBadRequest,
						errorJSON(app.MsgLength))
					app.Log.Error(app.MsgLength,
						"length", msgLength,
						"limit", app.MessageLimits.LengthChars,
						"user", req)
					return
				}

				msgCount := len(app.Messages)
				if msgCount >= app.MessageLimits.MaxCount {
					writeJSON(w, http.StatusBadRequest,
						errorJSON(app.MsgCount))
					app.Log.Error(app.MsgCount,
						"count", app.MessageLimits.MaxCount,
						"user", req)
					return
				}

				msgCount += 1
				message.Count = msgCount
				message.Data = formContent
				app.Messages = append(app.Messages, &message)
				app.Log.Info("added message",
					"count", msgCount,
					"length", msgLength,
					"user", req)
			}

			if req.IsBrowser {
				toRoot(w, r, app.Root)
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

		writeJSON(w, http.StatusOK, app.Messages)
	}
}
