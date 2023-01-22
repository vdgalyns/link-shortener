package main

import (
	"log"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/server"
)

func main() {
	cfg := config.NewConfig()
	srv := server.NewServer(cfg)
	log.Fatal(srv.ListenAndServe())
}
