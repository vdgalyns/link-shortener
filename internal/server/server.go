package server

import (
	"net/http"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/handler"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/router"
	"github.com/vdgalyns/link-shortener/internal/service"
)

func NewServer(config *config.Config) *http.Server {
	repositories := repository.NewRepository(config)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services, config)

	r := router.NewRouter(handlers)
	srv := &http.Server{Addr: config.ServerAddress, Handler: r}
	return srv
}
