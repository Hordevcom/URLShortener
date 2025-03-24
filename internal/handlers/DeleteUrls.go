package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func (h *ShortenHandler) DeleteUrls(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var urlIDs []string

	err = json.Unmarshal(body, &urlIDs)
	if err != nil {
		http.Error(w, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	URLsCh := make(chan string, len(urlIDs))
	deleteCh := make(chan string, len(urlIDs))

	var wg sync.WaitGroup

	wg.Add(2)
	go h.UpdateDeleteWorker(r.Context(), URLsCh, &wg)
	go h.DeleteWorker(r.Context(), deleteCh, &wg)

	for _, id := range urlIDs {
		URLsCh <- id
		deleteCh <- id
	}
	close(URLsCh)
	close(deleteCh)

	wg.Wait()

	w.WriteHeader(http.StatusAccepted)
}

func (h *ShortenHandler) UpdateDeleteWorker(ctx context.Context, urlsCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for urlID := range urlsCh {
		h.DB.UpdateDeleteParam(ctx, urlID)
	}
}

func (h *ShortenHandler) DeleteWorker(ctx context.Context, urlsCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for urlID := range urlsCh {
		h.DB.Delete(ctx, urlID)
	}
}
