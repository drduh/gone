package server

import (
	"log"
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
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("serve failed")
	}
	return nil
}
