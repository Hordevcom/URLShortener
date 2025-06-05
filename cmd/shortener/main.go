package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/handlers"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
	"github.com/Hordevcom/URLShortener/internal/workers"

	"github.com/Hordevcom/URLShortener/internal/routes"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

// Переменные для версии сборки
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	certFile := "server.crt"
	keyFile := "server.key"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	DeleteCh := make(chan string, 6)
	logger := logging.NewLogger()
	JSONStorage := storage.NewJSONStorage()
	conf := config.NewConfig()
	strg := storage.NewStorage(conf, logger)
	db := pg.NewPGDB(conf, logger)
	handler := handlers.NewShortenHandler(
		strg, conf, *JSONStorage, *db, DeleteCh)
	router := routes.NewRouter(*handler)
	workers := workers.NewDeleteWorker(ctx, db, DeleteCh, *handler)

	if conf.DatabaseDsn != "" {
		pg.InitMigrations(conf, logger)
	}

	server := &http.Server{Addr: conf.ServerAdress, Handler: router}

	if !conf.HttpsEnable {
		go func() {
			logger.Infow("Starting http server", "addr", conf.ServerAdress)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalw("create server error: ", err)
			}
		}()
	} else {
		go func() {
			logger.Infow("Starting https server")
			err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil)
			if err != nil {
				logger.Fatal("error while start server: ", err)
			}
		}()
	}

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorw("Server shutdown error", err)
	}

	workers.StopWorker()
	handler.CloseCh()
}
