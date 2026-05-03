package handlers

import (
	"net/http"
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
					"count", app.NumMessages, "user", req)
				app.ClearMessages()
				app.Log.Info("cleared messages", "user", req)
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

			formContent := r.FormValue(formFieldMessage)
			if formContent != "" {
				message.Count++
				message.Data = formContent
				app.Messages[message.Count] = &message
				app.Log.Info("added message", "user", req)
			}

			if req.IsBrowser {
				toRoot(w, r, app.Root)
				return
			}
		}

		if r.URL.Query().Get("download") == "all" {
			app.Log.Debug("downloading messages",
				"count", app.NumMessages, "user", req)
			app.ServeMessages(w)
			return
		}

		writeJSON(w, http.StatusOK, app.Messages)
	}
}
