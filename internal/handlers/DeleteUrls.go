package handlers

import (
	"encoding/json"
	"io"
	"net/http"
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

	for _, id := range urlIDs {
		h.DeleteCh <- id
	}

	w.WriteHeader(http.StatusAccepted)
}

// func (h *ShortenHandler) UpdateDeleteWorker(ctx context.Context, urlsCh <-chan string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for urlID := range urlsCh {
// 		h.DB.UpdateDeleteParam(ctx, urlID)
// 	}
// }

// func (h *ShortenHandler) DeleteWorker(ctx context.Context, urlsCh <-chan string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for urlID := range urlsCh {
// 		h.DB.Delete(ctx, urlID)
// 	}
// }
