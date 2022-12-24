package handler

import (
	"github.com/vdgalyns/link-shortener/internal/service"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	services *service.Service
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Add(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var _, id, _ = strings.Cut(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "url is incorrect", http.StatusBadRequest)
		return
	}
	v, err := h.services.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(b) == 0 {
		http.Error(w, "url not passed", http.StatusBadRequest)
		return
	}
	url := string(b)
	if _, domain, _ := strings.Cut(url, "//"); len(domain) == 0 {
		http.Error(w, "this is not an url", http.StatusBadRequest)
		return
	}
	id, err := h.services.Add(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output := "http://" + r.Host + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(output))
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
