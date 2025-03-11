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
		ip, ua := r.RemoteAddr, r.UserAgent()
		tmplName := "index"

		tmpl, err := template.New(tmplName).Parse(templates.IndexPage)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, responseErrorTmplParse)
			app.Log.Error(errorTmplParse,
				"template", tmplName,
				"ip", ip, "ua", ua)
			return
		}

		response := templates.Index{
			Title:      app.Version,
			PathUpload: app.Settings.Paths.Upload,
			AuthUpload: app.Settings.Auth.Require.Upload,
			AuthHeader: "X-Auth",
			AuthHolder: "AAAA-BBBB",
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
