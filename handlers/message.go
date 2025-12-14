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
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if r.Method == http.MethodPost {
			if r.FormValue("clear") != "" {
				app.Log.Debug("clearing messages",
					"count", app.CountMessages(), "user", req)
				app.ClearMessages()
				app.Log.Info("cleared messages", "user", req)
			}

			message := storage.Message{
				Count: app.CountMessages(),
				Owner: storage.Owner{
					Agent: req.Agent,
					Mask:  req.Mask,
				},
				Time: storage.Time{
					Allow: time.Now().Format(app.TimeFormat),
				},
			}

			content := r.FormValue("message")
			if content != "" {
				message.Count++
				message.Data = content
				app.Log.Debug("adding message",
					"count", message.Count,
					"content", message.Data, "user", req)
				app.Messages[message.Count] = &message
				app.Log.Info("added message", "user", req)
			}

			toRoot(w, r, app.Root)
		}

		if r.URL.Query().Get("download") == "all" {
			app.Log.Debug("serving all messages",
				"count", app.CountMessages(), "user", req)
			app.ServeMessages(w)
			return
		}

		writeJSON(w, http.StatusOK, app.Messages)
	}
}
