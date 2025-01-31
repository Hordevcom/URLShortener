package config

import "flag"

type Config struct {
	ServerAdress string
	Host         string
}

func GetCLParams() Config {
	var conf Config
	flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server adress")
	flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")

	flag.Parse()
	return conf
}
