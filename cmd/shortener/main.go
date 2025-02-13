package main

import (
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func main() {
	logger := logging.NewLogger()
	strg := storage.NewStorage()
	conf := config.NewConfig()
	app := app.NewApp(strg, conf)
	router := routes.NewRouter(*app, logger)

	logger.Infow("Starting server", "addr", conf.ServerAdress)
	http.ListenAndServe(conf.ServerAdress, router)
}
