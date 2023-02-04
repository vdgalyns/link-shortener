package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", ":8080", "server address")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base url")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&cfg.FileStoragePath, "f", "data.json", "file storage path")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:password@localhost:5432/my_db", "database dsn")
	}
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
