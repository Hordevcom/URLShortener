package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	Host         string `env:"BASE_URL"`
}

func NewConfig() Config {
	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.Host != "" && conf.ServerAdress != "" {
		return conf
	}

	flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server adress")
	flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")

	flag.Parse()
	return conf
}
