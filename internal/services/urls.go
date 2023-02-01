package services

import (
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"github.com/vdgalyns/link-shortener/internal/repositories"
)

type Urls struct {
	repositories *repositories.Repositories
	config       *config.Config
}

func (u *Urls) Get(id string) (entities.URL, error) {
	_, err := entities.ValidateUrlId(id)
	if err != nil {
		return entities.URL{}, err
	}
	return u.repositories.Get(id)
}

func (u *Urls) Add(originalUrl string) (string, error) {
	valid := entities.ValidateUrl(originalUrl)
	if !valid {
		return "", ErrUrlNotValid
	}
	id, err := entities.CreateUrlId()
	if err != nil {
		return "", err
	}
	url := entities.URL{
		ID:          id,
		OriginalURL: originalUrl,
	}
	err = u.repositories.Add(url)
	if err != nil {
		return "", err
	}
	readyUrl := u.config.BaseURL + "/" + id
	return readyUrl, nil
}

func (u *Urls) GetAllByUserId(userId string) ([]entities.URL, error) {
	_, err := entities.ValidateUserId(userId)
	if err != nil {
		return nil, err
	}
	return u.repositories.GetAllByUserId(userId)
}

func NewUrls(repositories *repositories.Repositories, config *config.Config) *Urls {
	return &Urls{
		repositories: repositories,
		config:       config,
	}
}
