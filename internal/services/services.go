package services

import (
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"github.com/vdgalyns/link-shortener/internal/repositories"
)

type Kind interface {
	Get(id string) (entities.URL, error)
	Add(originalUrl string) (string, error)
	GetAllByUserId(userId string) ([]entities.URL, error)
}

type Services struct {
	Kind
}

func NewServices(repositories *repositories.Repositories, config *config.Config) *Services {
	return &Services{
		Kind: NewUrls(repositories, config),
	}
}