package app

import (
	"github.com/vdgalyns/link-shortener/internal/handler"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/service"
	"log"
	"net/http"
)

func NewApp() {
	repositories := repository.NewRepository()
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	handlerIndex := http.HandlerFunc(handlers.Index)
	server := &http.Server{
		Addr: ":8080",
		Handler: handlerIndex,
	}
	log.Fatal(server.ListenAndServe())
}
