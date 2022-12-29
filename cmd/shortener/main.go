package main

import (
	"github.com/vdgalyns/link-shortener/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()
	log.Fatal(srv.ListenAndServe())
}
