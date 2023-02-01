package repositories

import (
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Kind interface {
	Get(hash string) (entities.URL, error)
	Add(url entities.URL) error
	GetAllByUserId(userId string) ([]entities.URL, error)
}

type Repositories struct {
	Kind
}

func NewRepositories(config *config.Config) *Repositories {
	repositories := &Repositories{
		Kind: NewFile(config.FileStoragePath),
	}
	if len(config.FileStoragePath) == 0 {
		repositories.Kind = NewMemory()
	}
	return repositories
}
