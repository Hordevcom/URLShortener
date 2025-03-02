package main

import (
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"

	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func main() {
	logger := logging.NewLogger()
	JSONStorage := storage.NewJSONStorage()
	conf := config.NewConfig()
	strg := storage.NewStorage()
	file := files.NewFile(conf, logger, strg)
	pg := pg.NewPGDB(conf, logger, strg)
	app := app.NewApp(strg, conf, *JSONStorage, *file, pg)
	router := routes.NewRouter(*app)

	logger.Infow("Starting server", "addr", conf.ServerAdress)
	err := http.ListenAndServe(conf.ServerAdress, router)

	if err != nil {
		logger.Fatalw("create server error:1 ", err)
	}
}
