package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

// Message handles requests to post, read, clear
// and download Messages from Storage.
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
				app.Messages[message.Count] = &message
				app.Log.Debug("added message",
					"count", message.Count,
					"content", message.Data, "user", req)
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
