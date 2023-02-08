package handlers

import (
	"encoding/json"
	"errors"
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
type BodyAddBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type ResponseAddBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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
	value, _, err := cookies.ReadAndCreateCookieUserID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := h.services.Add(string(body), value)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, services.ErrURLIsExist) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(shortURL))
}

func (h *Handlers) AddJSON(w http.ResponseWriter, r *http.Request) {
	value, _, err := cookies.ReadAndCreateCookieUserID(w, r)
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
	shortURL, err := h.services.Add(body.URL, value)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, services.ErrURLIsExist) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	output := ResponseAddJSON{Result: shortURL}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(outputJSON)
}

func (h *Handlers) GetUrls(w http.ResponseWriter, r *http.Request) {
	value, exist, err := cookies.ReadAndCreateCookieUserID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !exist {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}
	links, err := h.services.GetAllByUserID(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(links) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}
	output := make([]ResponseGetUrls, 0)
	for _, link := range links {
		responseURL := ResponseGetUrls{OriginalURL: link.OriginalURL, ShortURL: h.config.BaseURL + "/" + link.Hash}
		output = append(output, responseURL)
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(outputJSON)
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.services.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte{})
}

func (h *Handlers) AddBatch(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadSigned(r, "user_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body := make([]BodyAddBatch, 0)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	originalURLs := make([]string, 0, len(body))
	for _, v := range body {
		originalURLs = append(originalURLs, v.OriginalURL)
	}
	readyUrls, err := h.services.AddBatch(originalURLs, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output := make([]ResponseAddBatch, 0, len(body))
	for i, v := range body {
		outputItem := ResponseAddBatch{CorrelationID: v.CorrelationID, ShortURL: readyUrls[i]}
		output = append(output, outputItem)
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(outputJSON)
}

func NewHandlers(services *services.Services, config *config.Config) *Handlers {
	return &Handlers{
		services: services,
		config:   config,
	}
}
