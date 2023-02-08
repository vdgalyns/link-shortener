package services

import (
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"github.com/vdgalyns/link-shortener/internal/repositories"
)

type Kind interface {
	Get(hash string) (entities.Link, error)
	Add(originalURL, userID string) (string, error)
	GetAllByUserID(userID string) ([]entities.Link, error)
	Ping() error
	AddBatch(originalURLs []string, userID string) ([]string, error)
}

type Services struct {
	Kind
}

func NewServices(repositories *repositories.Repositories, config *config.Config) *Services {
	return &Services{
		Kind: NewUrls(repositories, config),
	}
}
