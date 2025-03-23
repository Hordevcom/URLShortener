package handlers

import (
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
	go h.UpdateDeleteWorker(URLsCh, &wg)
	go h.DeleteWorker(deleteCh, &wg)

	for _, id := range urlIDs {
		URLsCh <- id
		deleteCh <- id
	}
	close(URLsCh)
	close(deleteCh)

	wg.Wait()

	w.WriteHeader(http.StatusAccepted)
}

func (h *ShortenHandler) UpdateDeleteWorker(urlsCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for urlID := range urlsCh {
		h.Db.UpdateDeleteParam(urlID)
	}
}

func (h *ShortenHandler) DeleteWorker(urlsCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for urlID := range urlsCh {
		h.Db.Delete(urlID)
	}
}
