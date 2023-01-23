package repository

import "github.com/vdgalyns/link-shortener/internal/config"

type Link interface {
	Get(hash string) (string, error)
	Add(hash string, url string) error
}

type Repository struct {
	Link
}

func NewRepository(config *config.Config) *Repository {
	if len(config.FileStoragePath) == 0 {
		return &Repository{Link: NewLinkMemoryRepository()}
	}
	return &Repository{Link: NewLinkFileRepository(config.FileStoragePath)}
}
