package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Serves index page with app routing features
func Index(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ua := r.RemoteAddr, r.UserAgent()

		if r.Method == http.MethodPost {
			if app.Settings.Auth.Require.Message &&
				!auth.Basic(app.Settings.Auth.Header, app.Settings.Auth.Token, r) {
				writeJSON(w, http.StatusUnauthorized, responseErrorDeny)
				app.Log.Error(errorDeny,
					"action", "message",
					"ip", ip, "ua", ua)
				return
			}

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
			writeJSON(w, http.StatusInternalServerError, responseErrorTmplParse)
			app.Log.Error(errorTmplParse,
				"template", tmplName,
				"error", err.Error(),
				"ip", ip, "ua", ua)
			return
		}

		settings := app.Settings
		duration := settings.Limits.Expiration.Duration
		response := templates.Index{
			Auth:            settings.Auth,
			DefaultDuration: duration.String(),
			Paths:           settings.Paths,
			Storage:         app.Storage,
			Theme:           theme,
			ThemePick:       settings.Index.ThemePick,
			Title:           settings.Index.Title,
			Version:         app.Version,
			VersionFull:     app.VersionFull,
		}

		if err = tmpl.Execute(w, response); err != nil {
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
