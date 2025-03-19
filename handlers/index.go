package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Serves index page with app routing features
func Index(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if r.Method == http.MethodPost {
			if r.FormValue("clear") != "" {
				app.Storage.ClearMessages()
				app.Log.Debug("cleared messages",
					"ip", ip, "ua", ua)
			}

			message := config.Message{
				Count: app.Storage.CountMessages(),
				Owner: config.Owner{
					Address: ip,
					Agent:   ua,
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
					"message", content,
					"ip", ip, "ua", ua)
			}
		}

		tmplName := "index"
		tmpl, err := template.New(tmplName).Parse(templates.HtmlIndex)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, responseErrorTmplParse)
			app.Log.Error(errorTmplParse,
				"template", tmplName,
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}

		settings := app.Settings
		auth := settings.Auth
		index := settings.Index
		paths := settings.Paths
		duration := settings.Limits.Expiration.Duration

		response := templates.Index{
			AuthHeader:      auth.Header,
			AuthHolder:      auth.Holder,
			AuthDownload:    auth.Require.Download,
			AuthList:        auth.Require.List,
			AuthUpload:      auth.Require.Upload,
			Files:           app.Storage.Files,
			DefaultDuration: duration.String(),
			Messages:        app.Storage.Messages,
			PathDownload:    paths.Download,
			PathList:        paths.List,
			PathUpload:      paths.Upload,
			PathHeartbeat:   paths.Heartbeat,
			Style:           index.Style,
			Title:           index.Title,
			Version:         app.Version,
			VersionFull:     app.VersionFull,
		}

		err = tmpl.Execute(w, response)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, responseErrorTmplExec)
			app.Log.Error(errorTmplExec,
				"template", tmplName,
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}

		app.Log.Info("served index",
			"ip", ip, "ua", ua)
	}
}
