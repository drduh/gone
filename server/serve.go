package server

import (
	"net/http"

	"github.com/drduh/gone/handlers"
)

func Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Heartbeat())
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return srv.ListenAndServe()
}
