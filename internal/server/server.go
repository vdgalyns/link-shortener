package server

import (
	"database/sql"
	"net/http"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/handlers"
	"github.com/vdgalyns/link-shortener/internal/repositories"
	"github.com/vdgalyns/link-shortener/internal/router"
	"github.com/vdgalyns/link-shortener/internal/services"
)

func NewServer(config *config.Config, database *sql.DB) *http.Server {
	repo := repositories.NewRepositories(config, database)
	serv := services.NewServices(repo, config)
	hand := handlers.NewHandlers(serv, config)

	r := router.NewRouter(hand)
	srv := &http.Server{Addr: config.ServerAddress, Handler: r}
	return srv
}
