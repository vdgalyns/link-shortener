package main

import (
	"log"

	"github.com/vdgalyns/link-shortener/internal/server"
)

func main() {
	srv := server.NewServer()
	log.Fatal(srv.ListenAndServe())
}
