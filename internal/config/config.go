package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
)

// Config структура
type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS" json:"server_address"`
	Host         string `env:"BASE_URL" json:"base_url"`
	FilePath     string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDsn  string `env:"DATABASE_DSN" json:"database_dsn"`
	HTTPSEnable  bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	ConfigFile   string `env:"CONFIG"`
}

var once sync.Once

// Конструктор для конфига
func NewConfig() Config {

	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.Host != "" && conf.ServerAdress != "" {
		return conf
	}

	once.Do(func() {

		if conf.DatabaseDsn == "" {
			flag.StringVar(&conf.DatabaseDsn, "d", "", "database dsn") //"postgres://postgres:1@localhost:5432/postgres"
		}

		if conf.FilePath == "" {
			flag.StringVar(&conf.FilePath, "f", "", "path to file") //"storage.txt"
		}

		flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server adress") //localhost:8080
		flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")
		flag.StringVar(&conf.ConfigFile, "c", "", "config file")
		flag.BoolVar(&conf.HTTPSEnable, "s", false, "use https or not")

		flag.Parse()
	})

	if conf.DatabaseDsn == "" && conf.FilePath == "" && conf.Host == "" && conf.ServerAdress == "" && conf.ConfigFile != "" {
		conf, err := loadConfFromFile(conf.ConfigFile)
		if err != nil {
			fmt.Println("error! ", err)
		}
		return *conf
	}
	return conf
}

// loadConfFromFile сканирует конфигурационные данные из файла
func loadConfFromFile(path string) (*Config, error) {
	fmt.Println("Load config from file")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
