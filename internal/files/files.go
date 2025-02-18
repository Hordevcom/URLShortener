package files

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/Hordevcom/URLShortener/internal/config"
	"go.uber.org/zap"
)

type JSONStruct struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type File struct {
	config config.Config
	logger zap.SugaredLogger
}

func NewFile(config config.Config, logger zap.SugaredLogger) *File {
	return &File{config: config, logger: logger}
}

var UUID int = 0

func UpdateFile(jsonStruct JSONStruct) {

	UUID++
	jsonStruct.UUID = strconv.Itoa(UUID)
	file, err := os.OpenFile("storage.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		errors.New("Failed to create file")
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(jsonStruct)

	if err != nil {
		errors.New("Failed marshal")
		return
	}
	jsonData = append(jsonData, '\n')

	_, err = file.Write(jsonData)

	if err != nil {
		errors.New("Failed to write file")
		return
	}
}

func (f *File) ReadFile() map[string]string {
	var jsonStrct JSONStruct
	storage := make(map[string]string)
	file, err := os.OpenFile(f.config.FilePath, os.O_RDONLY|os.O_CREATE, 06666)

	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &jsonStrct)
		storage[jsonStrct.ShortURL] = jsonStrct.OriginalURL
	}

	UUID, _ = strconv.Atoi(jsonStrct.UUID)
	return storage
}
