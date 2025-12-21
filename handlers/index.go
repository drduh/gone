package handlers

import (
	"html/template"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/templates"
)

const templateData = "data/*.tmpl"

// Index handles requests to load and render the index page.
func Index(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authRequest(w, r, app)
		if req == nil {
			return
		}
		app.Log.Info("serving index", "user", req)

		theme := getDefaultTheme(app.Style.Theme)
		if app.Style.AllowPick {
			theme = getTheme(w, r, theme,
				app.Cookie.Id,
				app.Cookie.Time.GetDuration(),
				app.Style.Available)
		}

		tmpl, err := template.New("index").ParseFS(templates.All, templateData)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorJSON(app.TmplParse))
			app.Log.Error(app.TmplParse, "error", err.Error(), "user", req)
			return
		}

		app.UpdateTimeRemaining()
		response := templates.Index{
			Auth:            app.Auth,
			Default:         app.Default,
			DefaultDuration: app.Expiration.String(),
			Hostname:        app.Hostname,
			Index:           app.Index,
			Limit:           app.Limit,
			NoFiles:         app.NoFiles,
			Paths:           app.Paths,
			ShowBuild:       app.ShowBuild,
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
