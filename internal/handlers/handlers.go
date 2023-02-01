package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/cookies"
	"github.com/vdgalyns/link-shortener/internal/entities"
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
	value, err := cookies.ReadSigned(r, "user_id")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie), errors.Is(err, cookies.ErrInvalidValue):
			value, err = entities.CreateUserID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = cookies.WriteSigned(w, http.Cookie{Name: "user_id", Value: value})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusBadRequest)
			return
		}
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
	value, err := cookies.ReadSigned(r, "user_id")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie), errors.Is(err, cookies.ErrInvalidValue):
			value, err = entities.CreateUserID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = cookies.WriteSigned(w, http.Cookie{Name: "user_id", Value: value})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusBadRequest)
			return
		}
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
	outputJson, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(outputJson)
}

func (h *Handlers) GetUrls(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadSigned(r, "user_id")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie), errors.Is(err, cookies.ErrInvalidValue):
			value, err = entities.CreateUserID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = cookies.WriteSigned(w, http.Cookie{Name: "user_id", Value: value})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte{})
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusBadRequest)
		}
		return
	}
	urls, err := h.services.GetAllByUserId(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output := make([]ResponseGetUrls, 0)
	for _, v := range urls {
		output = append(output, ResponseGetUrls{OriginalURL: v.OriginalURL, ShortURL: h.config.BaseURL + "/" + v.Hash})
	}
	outputJson, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(outputJson)
}

func NewHandlers(services *services.Services, config *config.Config) *Handlers {
	return &Handlers{
		services: services,
		config:   config,
	}
}
