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

		theme := getTheme(app.Theme)

		if app.ThemePick {
			cookieDuration := app.Cookie.Time.GetDuration()
			cookieNew := &http.Cookie{
				Name:    app.Cookie.Id,
				Expires: time.Now().Add(cookieDuration),
				Path:    "/",
			}

			themeForm := r.FormValue("theme")
			if themeForm != "" {
				theme = themeForm
				cookieNew.Value = theme
				http.SetCookie(w, cookieNew)
			} else {
				cookie, err := r.Cookie(app.Cookie.Id)
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
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.TmplParse))
			app.Log.Error(app.TmplParse,
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
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.TmplExec))
			app.Log.Error(app.TmplExec,
				"template", tmplName, "error", err.Error(), "user", req)
			return
		}
	}
}
