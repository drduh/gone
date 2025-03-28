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
		req := parseRequest(r)

		if r.Method == http.MethodPost {
			if !isAllowed(app, r) {
				deny(w, app, req)
				return
			}

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
		}

		theme := getTheme(app.Settings.Index.Theme)

		if app.Settings.Index.ThemePick {
			cookieDuration := app.Settings.Index.Cookie.Time.GetDuration()
			cookieNew := &http.Cookie{
				Name:    app.Settings.Index.Cookie.Id,
				Expires: time.Now().Add(cookieDuration),
				Path:    "/",
			}

			themeForm := r.FormValue("theme")
			if themeForm != "" {
				theme = themeForm
				cookieNew.Value = theme
				http.SetCookie(w, cookieNew)
			} else {
				cookie, err := r.Cookie(app.Settings.Index.Cookie.Id)
				if err != nil || cookie.Value == "" {
					cookieNew.Value = theme
					http.SetCookie(w, cookieNew)
				} else {
					theme = cookie.Value
				}
			}
		}

		tmplName := "index.tmpl"
		tmpl, err := template.New(tmplName).ParseFS(templates.All, "data/*.tmpl")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Error.TmplParse))
			app.Log.Error(app.Error.TmplParse,
				"template", tmplName, "error", err.Error(), "user", req)
			return
		}

		settings := app.Settings
		duration := settings.Limits.Expiration.Duration
		response := templates.Index{
			Auth:            settings.Auth,
			DefaultDuration: duration.String(),
			Index:           settings.Index,
			Limits:          app.Limits,
			Paths:           settings.Paths,
			Storage:         app.Storage,
			Theme:           theme,
			ThemePick:       settings.Index.ThemePick,
			Title:           settings.Index.Title,
			Version:         app.Version,
			VersionFull:     app.VersionFull,
		}

		if err = tmpl.Execute(w, response); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Error.TmplExec))
			app.Log.Error(app.Error.TmplExec,
				"template", tmplName, "error", err.Error(), "user", req)
			return
		}

		app.Log.Info("served index", "user", req)
	}
}
