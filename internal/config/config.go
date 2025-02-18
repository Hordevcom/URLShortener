package config

import (
	"flag"
	"fmt"

	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	Host         string `env:"BASE_URL"`
	FilePath     string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() Config {
	logger := logging.NewLogger()
	var conf Config
	err := env.Parse(&conf)

	logger.Infow(conf.FilePath)

	if err != nil {
		fmt.Println(err)
	}

	if conf.Host != "" && conf.ServerAdress != "" && conf.FilePath != "" {
		return conf
	}

	flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server adress")
	flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")
	flag.StringVar(&conf.FilePath, "f", "storage.txt", "path to file")

	flag.Parse()
	return conf
}
