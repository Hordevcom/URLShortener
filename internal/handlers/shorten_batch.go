package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
)

type ShortenRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenResponce struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *ShortenHandler) BatchShortenURL(w http.ResponseWriter, r *http.Request) {
	var requests []ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&requests)

	if err != nil {
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	if len(requests) == 0 {
		http.Error(w, "Batch cannot be empty", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.DB.Begin(r.Context())
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(r.Context())

	query := `INSERT INTO urls (short_url, original_url, user_id)
	 VALUES ($1, $2, $3) ON CONFLICT (short_url) DO NOTHING`
	var responces []ShortenResponce
	for _, req := range requests {
		shortURL := fmt.Sprintf("%x", md5.Sum([]byte(req.OriginalURL)))[:8]

		_, err := tx.Exec(r.Context(), query, shortURL, req.OriginalURL, 0)

		if err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		responces = append(responces, ShortenResponce{
			CorrelationID: req.CorrelationID,
			ShortURL:      h.Config.Host + "/" + shortURL,
		})

	}

	if err := tx.Commit(r.Context()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responces)
}
