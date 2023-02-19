package repositories

import (
	"database/sql"

	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Kind interface {
	Get(hash string) (entities.Link, error)
	GetByOriginalURL(originalURL string) (entities.Link, error)
	Add(link entities.Link) error
	GetAllByUserID(userID string) ([]entities.Link, error)
	Ping() error
	AddBatch(links []entities.Link) error
	RemoveBatch(urlHashes []string, userID string) error
}

type Repositories struct {
	Kind
}

func NewRepositories(config *config.Config, database *sql.DB) *Repositories {
	repositories := new(Repositories)
	if len(config.DatabaseDSN) != 0 {
		repositories.Kind = NewDatabase(database)
	} else {
		if len(config.FileStoragePath) != 0 {
			repositories.Kind = NewFile(config.FileStoragePath)
		} else {
			repositories.Kind = NewMemory()
		}
	}
	return repositories
}
