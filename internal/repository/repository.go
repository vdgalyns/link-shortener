package repository

import (
	"github.com/vdgalyns/link-shortener/internal/config"
)

type Linker interface {
	Get(hash string) (string, error)
	Add(hash string, url string) error
}

type Repository struct {
	Linker
}

func NewRepository(config *config.Config) *Repository {
	if len(config.FileStoragePath) == 0 {
		return &Repository{NewLinkMemoryRepository()}
	}
	return &Repository{NewLinkFileRepository(config.FileStoragePath)}
}
