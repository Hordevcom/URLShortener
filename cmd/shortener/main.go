package main

import (
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func main() {
	strg := storage.NewStorage()
	conf := config.NewConfig()
	app := app.NewApp(strg, conf)
	router := routes.NewRouter(*app)
	http.ListenAndServe(conf.ServerAdress, router)
}
