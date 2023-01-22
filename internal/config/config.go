package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func NewConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
