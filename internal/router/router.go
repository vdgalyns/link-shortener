package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vdgalyns/link-shortener/internal/handler"
	middle "github.com/vdgalyns/link-shortener/internal/router/middleware"
)

func NewRouter(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middle.GzipHandle)

	r.Get("/{hash}", h.Get)
	r.Post("/", h.Add)
	r.Post("/api/shorten", h.AddWithJSON)

	return r
}
