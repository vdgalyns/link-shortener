package main

import (
	"log"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	srv := server.NewServer(cfg)
	log.Fatal(srv.ListenAndServe())
}
