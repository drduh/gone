package handlers

import (
	"net/http"
	"time"

	"github.com/drduh/gone/config"
)

// Handle text messages
func Message(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, allowed := authRequest(w, r, app)
		if !allowed {
			return
		}

		if r.Method == http.MethodPost {
			if r.FormValue("clear") != "" {
				app.ClearMessages()
				app.Log.Debug("cleared messages", "user", req)
			}

			message := config.Message{
				Count: app.CountMessages(),
				Owner: config.Owner{
					Address: req.Address,
					Agent:   req.Agent,
				},
				Time: config.Time{
					Allow: time.Now().Format(app.TimeFormat),
				},
			}

			content := r.FormValue("message")
			if content != "" {
				message.Count++
				message.Data = content
				app.Messages[message.Count] = &message
				app.Log.Debug("added message",
					"message", content, "user", req)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		if r.URL.Query().Get("download") == "all" {
			app.Storage.ServeMessages(w)
			return
		}

		writeJSON(w, http.StatusOK, app.Messages)
	}
}
