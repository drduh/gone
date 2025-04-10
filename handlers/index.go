package handlers

import (
	"html/template"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

// Index handles requests for the main index page
// with all available application features.
func Index(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseRequest(r)

		if !app.Allow(app.PerMinute) {
			writeJSON(w, http.StatusTooManyRequests, errorJSON(app.RateLimit))
			app.Log.Error(app.RateLimit, "user", req)
			return
		}

		app.Log.Info("serving index", "user", req)

		theme := getDefaultTheme(app.Style.Theme)
		app.Log.Debug("got theme", "default", theme)
		if app.Style.AllowPick {
			theme = getTheme(w, r, theme,
				app.Cookie.Id, app.Cookie.Time.GetDuration())
			app.Log.Debug("got theme", "selected", theme)
		}

		tmpl, err := template.New("index").ParseFS(templates.All, "data/*.tmpl")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.TmplParse))
			app.Log.Error(app.TmplParse, "error", err.Error(), "user", req)
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
			Uptime:          app.Uptime(),
			Version:         app.Version,
		}

		if err = tmpl.Execute(w, response); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.TmplExec))
			app.Log.Error(app.TmplExec, "error", err.Error(), "user", req)
			return
		}
	}
}
