package server

import (
	"fmt"
	"net/http"

	"github.com/drduh/gone/config"
)

// Start HTTP server with configured paths
func Serve(app *config.App) error {
	go expiryWorker(app)
	app.Log.Info("starting server", "port", app.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: getHandler(app),
	}
	return server.ListenAndServe()
}
