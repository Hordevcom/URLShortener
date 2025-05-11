package files

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Hordevcom/URLShortener/internal/config"
	"go.uber.org/zap"
)

// JSONStruct используется для сериализации JSON
type JSONStruct struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// File структура
type File struct {
	config config.Config
	logger zap.SugaredLogger
	UUID   int
}

// NewFile конструктор для структуры File
func NewFile(config config.Config, logger zap.SugaredLogger) *File {
	f := &File{config: config, logger: logger}
	return f
}

// Set сохранение данных в файл
func (f *File) Set(ctx context.Context, shortURL, origURL string, userID int) bool { //jsonStruct JSONStruct,

	jsonStruct := JSONStruct{
		ShortURL:    shortURL,
		OriginalURL: origURL,
	}

	f.UUID++
	jsonStruct.UUID = strconv.Itoa(f.UUID)

	err := os.MkdirAll(filepath.Dir(f.config.FilePath), os.ModePerm)

	if err != nil {
		return false
	}

	file, err := os.OpenFile(f.config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return false
	}
	defer file.Close()

	jsonData, err := json.Marshal(jsonStruct)

	if err != nil {
		f.logger.Errorw("Failed marshal")
		return false
	}
	jsonData = append(jsonData, '\n')

	_, err = file.Write(jsonData)

	if err != nil {
		f.logger.Errorw("Failed to write file")
		return false
	}

	return true
}

// Get получить данные из файла
func (f *File) Get(ctx context.Context, shortURL string) (string, bool) { //strg storage.Storage

	data := make(map[string]string)

	var jsonStrct JSONStruct
	err := os.MkdirAll(filepath.Dir(f.config.FilePath), os.ModePerm)

	if err != nil {
		return "", false
	}

	file, err := os.OpenFile(f.config.FilePath, os.O_RDONLY|os.O_CREATE, 06666)

	if err != nil {
		return "", false
	}

	f.logger.Infow("created file in direction: " + f.config.FilePath)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &jsonStrct)
		data[jsonStrct.ShortURL] = jsonStrct.OriginalURL
	}

	f.UUID, _ = strconv.Atoi(jsonStrct.UUID)

	origURL, found := data[shortURL]
	if !found {
		return "", false
	}

	return origURL, true
}
