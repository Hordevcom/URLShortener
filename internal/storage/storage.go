package storage

import (
	"context"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
	"go.uber.org/zap"
)

// Storage интерфейс для сохранения / загрузки данных
type Storage interface {
	Set(ctx context.Context, key, value string, userID int) bool
	Get(ctx context.Context, key string) (string, bool)
}

// MapStorage структура инициализации временного хранилища урлов
type MapStorage struct {
	data map[string]string
}

// NewMapStorage конструктор для MapStorage
func NewMapStorage() *MapStorage {
	return &MapStorage{data: make(map[string]string)}
}

// NewStorage инициализация storage
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

// Set сохранение данных в storage
func (s *MapStorage) Set(ctx context.Context, key, value string, userID int) bool {
	s.data[key] = value
	return true
}

// Get загрузка данных из storage
func (s *MapStorage) Get(ctx context.Context, key string) (string, bool) {
	value, exist := s.data[key]
	return value, exist
}

// JSONStorage структура для серриализации данных
type JSONStorage struct {
	URL string `json:"url"`
}

// NewJSONStorage конструктор для JSONStorage
func NewJSONStorage() *JSONStorage {
	return &JSONStorage{URL: ""}
}

// Get загрузка данных из NewJSONStorage
func (s *JSONStorage) Set(key, value string) {
	s.URL = value
}

// Set сохранение данных в NewJSONStorage
func (s *JSONStorage) Get() string {
	return s.URL
}
