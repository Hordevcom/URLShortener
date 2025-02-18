package files

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"go.uber.org/zap"
)

type JSONStruct struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type File struct {
	config  config.Config
	logger  zap.SugaredLogger
	storage storage.MapStorage
}

func NewFile(config config.Config, logger zap.SugaredLogger, storage storage.Storage) *File {
	f := &File{config: config, logger: logger}
	f.ReadFile(storage)
	return f
}

var UUID int = 0

func (f *File) UpdateFile(jsonStruct JSONStruct) {

	UUID++
	jsonStruct.UUID = strconv.Itoa(UUID)
	file, err := os.OpenFile("storage.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		f.logger.Errorw("Failed to create file")
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(jsonStruct)

	if err != nil {
		f.logger.Errorw("Failed marshal")
		return
	}
	jsonData = append(jsonData, '\n')

	_, err = file.Write(jsonData)

	if err != nil {
		f.logger.Errorw("Failed to write file")
		return
	}
}

func (f *File) ReadFile(strg storage.Storage) {
	var jsonStrct JSONStruct
	file, err := os.OpenFile(f.config.FilePath, os.O_RDONLY|os.O_CREATE, 06666)

	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &jsonStrct)
		strg.Set(jsonStrct.ShortURL, jsonStrct.OriginalURL)
	}

	UUID, _ = strconv.Atoi(jsonStrct.UUID)
}
