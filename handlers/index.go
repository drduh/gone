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
				app.Storage.Messages = make(map[int]*config.Message)
				app.Log.Debug("cleared messages",
					"ip", ip, "ua", ua)
			}

			message := config.Message{
				Count: len(app.Storage.Messages),
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
				"ip", ip, "ua", ua)
			return
		}

		response := templates.Index{
			Title:         "gone",
			Version:       app.Version,
			AuthHeader:    app.Settings.Auth.Header,
			AuthHolder:    app.Settings.Auth.Holder,
			AuthDownload:  app.Settings.Auth.Require.Download,
			AuthList:      app.Settings.Auth.Require.List,
			AuthUpload:    app.Settings.Auth.Require.Upload,
			PathDownload:  app.Settings.Paths.Download,
			PathList:      app.Settings.Paths.List,
			PathUpload:    app.Settings.Paths.Upload,
			PathHeartbeat: app.Settings.Paths.Heartbeat,
			Messages:      app.Storage.Messages,
		}

		err = tmpl.Execute(w, response)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, responseErrorTmplExec)
			app.Log.Error(errorTmplExec,
				"template", tmplName,
				"ip", ip, "ua", ua)
			return
		}

		app.Log.Info("served index",
			"ip", ip, "ua", ua)
	}
}
