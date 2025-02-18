package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	Host         string `env:"BASE_URL"`
	FilePath     string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() Config {

	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.FilePath == "" {
		flag.StringVar(&conf.FilePath, "f", "storage.txt", "path to file")
	}

	if conf.Host == "" {
		flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")
	}

	if conf.ServerAdress == "" {
		flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server adress")
	}

	flag.Parse()
	return conf
}
