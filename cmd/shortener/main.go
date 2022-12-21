package main

import (
	"github.com/vdgalyns/link-shortener/internal/app/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.Index)
	http.ListenAndServe(":8080", nil)
}
