package handlers

import (
	"github.com/vdgalyns/link-shortener/internal/app/datastore"
	"github.com/vdgalyns/link-shortener/internal/app/generator"
	"io"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		retrieve(w, r)
	case http.MethodPost:
		shorten(w, r)
	default:
		notfound(w, r)
	}
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func shorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(b) == 0 {
		http.Error(w, "Body is not empty", http.StatusBadRequest)
		return
	}
	url := string(b)
	if !strings.Contains(url, "http") {
		http.Error(w, "Body is not a link", http.StatusBadRequest)
		return
	}
	id := generator.Make()
	datastore.Datastore[id] = url
	protocol := "http://"
	output := protocol + r.Host + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(output))
}
func retrieve(w http.ResponseWriter, r *http.Request) {
	var _, id, _ = strings.Cut(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Link is not passed", http.StatusBadRequest)
		return
	}
	v, ok := datastore.Datastore[id]
	if !ok {
		http.Error(w, "Link is not exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}
