package server

import (
	"github.com/vdgalyns/link-shortener/internal/handler"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/router"
	"github.com/vdgalyns/link-shortener/internal/service"
	"net/http"
)

func NewServer() *http.Server {
	repositories := repository.NewRepository()
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	r := router.NewRouter(handlers)
	srv := &http.Server{Addr: ":8080", Handler: r}
	return srv
}
