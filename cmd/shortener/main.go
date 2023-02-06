package main

import (
	"fmt"
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
	db, err := database.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		// TODO: Сейчас игнонорирую из-за ручки /ping
	}
	srv := server.NewServer(cfg, db)
	log.Fatal(srv.ListenAndServe())
}
