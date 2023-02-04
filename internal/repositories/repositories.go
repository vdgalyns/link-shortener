package repositories

import (
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Kind interface {
	Get(hash string) (entities.URL, error)
	Add(url entities.URL) error
	GetAllByUserID(userID string) ([]entities.URL, error)
}
type KindDatabase interface {
	Kind
	Ping() error
}

type Repositories struct {
	Kind
	Database KindDatabase
}

func NewRepositories(config *config.Config) *Repositories {
	repositories := &Repositories{
		Kind:     NewFile(config.FileStoragePath),
		Database: NewDatabase(config.DatabaseDSN),
	}
	if len(config.FileStoragePath) == 0 {
		repositories.Kind = NewMemory()
	}
	return repositories
}
