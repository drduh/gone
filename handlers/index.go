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
		app.Log.Info("serving index", "user", req)

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

		response := templates.Index{
			Auth:            app.Auth,
			DefaultDuration: app.Expiration.String(),
			Hostname:        app.Hostname,
			Index:           app.Index,
			Limits:          app.Limits,
			Paths:           app.Paths,
			Storage:         app.Storage,
			Theme:           theme,
			ThemePick:       app.ThemePick,
			Uptime:          app.Uptime(),
			Version:         app.Version,
		}

		if err = tmpl.Execute(w, response); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.Error.TmplExec))
			app.Log.Error(app.Error.TmplExec,
				"template", tmplName, "error", err.Error(), "user", req)
			return
		}
	}
}
