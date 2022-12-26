package handler

import (
	"errors"
	"fmt"
	"github.com/vdgalyns/link-shortener/internal/service"
	"io"
	"net/http"
	"strings"
)

var (
	ErrReadLink    = errors.New("error read link")
	ErrLinkIsEmpty = errors.New("link cannot be empty")
)

type Handler struct {
	services *service.Service
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
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
	_, hash, _ := strings.Cut(r.URL.Path, "/")
	url, err := h.services.Link.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrReadLink.Error(), http.StatusBadRequest)
		return
	}
	if len(b) == 0 {
		http.Error(w, ErrLinkIsEmpty.Error(), http.StatusBadRequest)
		return
	}
	s := string(b)
	hash, err := h.services.Link.Add(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: статичный протокол (поменять на определение текущего)
	link := fmt.Sprintf("http://%s/%s", r.Host, hash)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(link))
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
