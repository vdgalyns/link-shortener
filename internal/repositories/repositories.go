package repositories

import (
	"database/sql"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Kind interface {
	Get(hash string) (entities.URL, error)
	Add(url entities.URL) error
	GetAllByUserID(userID string) ([]entities.URL, error)
	Ping() error
}
type KindDatabase interface {
	Kind
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
