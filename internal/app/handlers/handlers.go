package handlers

import (
	"github.com/vdgalyns/link-shortener/internal/app/generator"
	"github.com/vdgalyns/link-shortener/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		read(w, r)
	case http.MethodPost:
		add(w, r)
	default:
		notfound(w, r)
	}
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func add(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(b) == 0 {
		http.Error(w, "body is not empty", http.StatusBadRequest)
		return
	}
	url := string(b)
	if !strings.Contains(url, "http") {
		http.Error(w, "body is not a link", http.StatusBadRequest)
		return
	}
	id := generator.CreateRandomId()
	storage.Add(id, url)
	protocol := "http://"
	output := protocol + r.Host + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(output))
}
func read(w http.ResponseWriter, r *http.Request) {
	var _, id, _ = strings.Cut(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "link is not passed", http.StatusBadRequest)
		return
	}
	v, ok := storage.Get(id)
	if !ok {
		http.Error(w, "link is not exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}
