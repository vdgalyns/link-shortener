package main

import (
	"log"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/database"
	"github.com/vdgalyns/link-shortener/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, _ := database.NewDatabase(cfg.DatabaseDSN)
	srv := server.NewServer(cfg, db)
	log.Fatal(srv.ListenAndServe())
}
