package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Handle text messages
func Message(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if !isAllowed(app, r) {
			deny(w, app, req)
			return
		}

		if r.Method == http.MethodPost {
			if r.FormValue("clear") != "" {
				app.Storage.ClearMessages()
				app.Log.Debug("cleared messages", "user", req)
			}

			message := config.Message{
				Count: app.Storage.CountMessages(),
				Owner: config.Owner{
					Address: req.Address,
					Agent:   req.Agent,
				},
				Time: config.Time{
					Allow: time.Now().Format(app.Settings.Audit.TimeFormat),
				},
			}

			content := r.FormValue("message")
			if content != "" {
				message.Count++
				message.Data = content
				app.Storage.Messages[message.Count] = &message
				app.Log.Debug("added message",
					"message", content, "user", req)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		writeJSON(w, http.StatusOK, app.Storage.Messages)
	}
}
