package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *ShortenHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if originalURL, exists := h.Storage.Get(r.Context(), shortURL); exists {
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "URL not found", http.StatusGone)
	}
}
