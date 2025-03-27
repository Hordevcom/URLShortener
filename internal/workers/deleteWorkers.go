package workers

import (
	"context"
	"sync"

	"github.com/Hordevcom/URLShortener/internal/handlers"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
)

type Worker struct {
	deleteCh chan string
	DB       pg.PGDB
	ctx      context.Context
	wg       sync.WaitGroup
	handler  handlers.ShortenHandler
}

func NewDeleteWorker(ctx context.Context, DB *pg.PGDB, deleteCh chan string, handler handlers.ShortenHandler) *Worker {
	worker := &Worker{
		ctx:      ctx,
		DB:       *DB,
		deleteCh: deleteCh,
		handler:  handler,
	}

	worker.wg.Add(2)
	go worker.UpdateDeleteWorker(ctx, deleteCh)
	go worker.DeleteWorker(ctx, deleteCh)

	return worker
}

func (w *Worker) UpdateDeleteWorker(ctx context.Context, urlsCh <-chan string) {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case urlID, ok := <-w.deleteCh:
			if !ok {
				return
			}
			w.DB.UpdateDeleteParam(ctx, urlID)
			w.handler.AddToChan(urlID)
		}
	}
}

func (w *Worker) DeleteWorker(ctx context.Context, urlsCh <-chan string) {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case urlID, ok := <-w.deleteCh:
			if !ok {
				return
			}
			w.DB.Delete(ctx, urlID)
		}
	}
}

func (w *Worker) StopWorker() {
	w.wg.Wait()
}
