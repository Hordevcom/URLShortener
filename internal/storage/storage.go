package storage

import (
	"context"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
	"go.uber.org/zap"
)

type Storage interface {
	Set(ctx context.Context, key, value string, userID int) bool
	Get(ctx context.Context, key string) (string, bool)
}

type MapStorage struct {
	data map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{data: make(map[string]string)}
}

func NewStorage(conf config.Config, logger zap.SugaredLogger) Storage {
	if conf.DatabaseDsn != "" {
		logger.Infow("DB config")
		return pg.NewPGDB(conf, logger)
	}
	if conf.FilePath != "" {
		logger.Infow("file config")
		return files.NewFile(conf, logger)
	}
	logger.Infow("memory config")
	return NewMapStorage()
}

func (s *MapStorage) Set(ctx context.Context, key, value string, userID int) bool {
	s.data[key] = value
	return true
}

func (s *MapStorage) Get(ctx context.Context, key string) (string, bool) {
	value, exist := s.data[key]
	return value, exist
}

type JSONStorage struct {
	URL string `json:"url"`
}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{URL: ""}
}

func (s *JSONStorage) Set(key, value string) {
	s.URL = value
}

func (s *JSONStorage) Get() string {
	return s.URL
}
