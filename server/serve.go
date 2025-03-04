package server

import (
	"fmt"
	"net/http"

	"github.com/drduh/gone/config"
	"github.com/drduh/gone/handlers"
)

func Serve(c *config.App) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Heartbeat())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: mux,
	}
	return srv.ListenAndServe()
}
