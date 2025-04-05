package handlers

import (
	"html/template"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Serves index page with app routing features
func Index(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)
		app.Log.Info("serving index", "user", req)

		theme := getDefaultTheme(app.Theme)
		if app.ThemePick {
			theme = getTheme(w, r, theme,
				app.Cookie.Id, app.Cookie.Time.GetDuration())
		}

		tmplName := "index"
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
