package main

import (
	"context"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/handlers"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
	"github.com/Hordevcom/URLShortener/internal/workers"

	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func main() {

	DeleteCh := make(chan string, 6)
	logger := logging.NewLogger()
	JSONStorage := storage.NewJSONStorage()
	conf := config.NewConfig()
	strg := storage.NewStorage(conf, logger)
	db := pg.NewPGDB(conf, logger)
	handler := handlers.NewShortenHandler(
		strg, conf, *JSONStorage, *db, DeleteCh)
	router := routes.NewRouter(*handler)
	workers.NewDeleteWorker(context.Background(), db, DeleteCh)

	if conf.DatabaseDsn != "" {
		pg.InitMigrations(conf, logger)
	}

	logger.Infow("Starting server", "addr", conf.ServerAdress)
	err := http.ListenAndServe(conf.ServerAdress, router)

	if err != nil {
		logger.Fatalw("create server error: ", err)
	}

}
