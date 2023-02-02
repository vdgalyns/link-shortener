package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/cookies"
	"github.com/vdgalyns/link-shortener/internal/services"
)

type ResponseAddJSON struct {
	Result string `json:"result"`
}
type ResponseGetUrls struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Handlers struct {
	services *services.Services
	config   *config.Config
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, err := h.services.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	value, _, err := cookies.ReadAndCreateCookieUserId(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url, err := h.services.Add(string(body), value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}

func (h *Handlers) AddJSON(w http.ResponseWriter, r *http.Request) {
	value, _, err := cookies.ReadAndCreateCookieUserId(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body := struct {
		URL string `json:"url"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url, err := h.services.Add(body.URL, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output := ResponseAddJSON{Result: url}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(outputJSON)
}

func (h *Handlers) GetUrls(w http.ResponseWriter, r *http.Request) {
	value, exist, err := cookies.ReadAndCreateCookieUserId(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !exist {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}
	urls, err := h.services.GetAllByUserID(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output := make([]ResponseGetUrls, 0)
	for _, v := range urls {
		output = append(output, ResponseGetUrls{OriginalURL: v.OriginalURL, ShortURL: h.config.BaseURL + "/" + v.Hash})
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(outputJSON)
}

func NewHandlers(services *services.Services, config *config.Config) *Handlers {
	return &Handlers{
		services: services,
		config:   config,
	}
}
