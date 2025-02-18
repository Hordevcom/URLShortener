package main

import (
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func main() {
	logger := logging.NewLogger()
	JSONStorage := storage.NewJSONStorage()
	conf := config.NewConfig()
	file := files.NewFile(conf, logger)
	strg := storage.NewStorage(*file)
	app := app.NewApp(strg, conf, *JSONStorage)
	router := routes.NewRouter(*app)

	logger.Infow("Starting server", "addr", conf.ServerAdress)
	http.ListenAndServe(conf.ServerAdress, router)
}
