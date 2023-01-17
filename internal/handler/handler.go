package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vdgalyns/link-shortener/internal/service"
	"io"
	"net/http"
)

var (
	ErrReadLink    = errors.New("error read link")
	ErrLinkIsEmpty = errors.New("link cannot be empty")
)

type Handler struct {
	services *service.Service
}

type BodyWhenAdding struct {
	URL string `json:"url"`
}
type ResponseWhenAdding struct {
	Result string `json:"result"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
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
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(link))
}

func (h *Handler) AddWithJSON(w http.ResponseWriter, r *http.Request) {
	body := BodyWhenAdding{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body.URL) == 0 {
		http.Error(w, ErrLinkIsEmpty.Error(), http.StatusBadRequest)
		return
	}
	hash, err := h.services.Link.Add(body.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: статичный протокол (поменять на определение текущего)
	link := fmt.Sprintf("http://%s/%s", r.Host, hash)
	response := ResponseWhenAdding{
		Result: link,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
