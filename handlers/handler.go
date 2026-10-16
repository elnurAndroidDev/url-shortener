package handlers

import (
	"fmt"
	"net/http"
	"url-shortener/models"
	"url-shortener/storage"
)

type Handler struct {
	store *storage.Store
}

func NewHandler(store *storage.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
		return
	}

	shortURL := models.GenerateShortURL()
	h.store.Save(shortURL, longURL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://shorty.url/" + shortURL))
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	shortURL := r.URL.Path[1:]
	fmt.Println("shortURL: ", shortURL)

	longURL, exists := h.store.Get(shortURL)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
