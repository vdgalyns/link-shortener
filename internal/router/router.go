package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vdgalyns/link-shortener/internal/handlers"
	middle "github.com/vdgalyns/link-shortener/internal/router/middleware"
)

func NewRouter(h *handlers.Handlers) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middle.GzipDecompress)
	r.Use(middle.GzipCompress)
	// r.Use(middle.ReadAndWriteCookieUserID)

	r.Get("/{id}", h.Get)
	r.Post("/", h.Add)
	r.Post("/api/shorten", h.AddJSON)
	r.Get("/api/user/urls", h.GetUrls)
	r.Post("/api/shorten/batch", h.AddBatch)
	r.Get("/ping", h.Ping)
	r.Delete("/api/user/urls", h.DeleteBatch)

	return r
}
