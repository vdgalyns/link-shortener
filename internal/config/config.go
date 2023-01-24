package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

const (
	ServerAddress   string = ":8080"
	BaseURL         string = "http://localhost:8080"
	FileStoragePath string = "data.txt"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	a := flag.String("a", "", "server address")
	b := flag.String("b", "", "base url")
	f := flag.String("f", "", "file storage path")
	flag.Parse()
	if cfg.ServerAddress == "" {
		if *a == "" {
			cfg.ServerAddress = ServerAddress
		} else {
			cfg.ServerAddress = *a
		}
	}
	if cfg.BaseURL == "" {
		if *b == "" {
			cfg.BaseURL = BaseURL
		} else {
			cfg.BaseURL = *b
		}
	}
	if cfg.FileStoragePath == "" {
		if *f == "" {
			cfg.FileStoragePath = FileStoragePath
		} else {
			cfg.FileStoragePath = *f
		}
	}
	return &cfg, nil
}
