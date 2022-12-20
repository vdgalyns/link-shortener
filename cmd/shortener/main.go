package main

import (
	"github.com/vdgalyns/link-shortener/internal/app/handlers"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(handlers.Handler)
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	server.ListenAndServe()
}
