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

func (u *Urls) Get(hash string) (entities.URL, error) {
	_, err := entities.ValidateURLHash(hash)
	if err != nil {
		return entities.URL{}, err
	}
	return u.repositories.Get(hash)
}

func (u *Urls) Add(originalURL, userID string) (string, error) {
	valid := entities.ValidateURL(originalURL)
	if !valid {
		return "", ErrURLNotValid
	}
	hash, err := entities.CreateURLHash(originalURL)
	if err != nil {
		return "", err
	}
	url := entities.URL{
		Hash:        hash,
		UserID:      userID,
		OriginalURL: originalURL,
	}
	err = u.repositories.Add(url)
	if err != nil {
		return "", err
	}
	readyUrl := u.config.BaseURL + "/" + hash
	return readyUrl, nil
}

func (u *Urls) GetAllByUserID(userID string) ([]entities.URL, error) {
	_, err := entities.ValidateUserID(userID)
	if err != nil {
		return nil, err
	}
	return u.repositories.GetAllByUserID(userID)
}

func NewUrls(repositories *repositories.Repositories, config *config.Config) *Urls {
	return &Urls{
		repositories: repositories,
		config:       config,
	}
}
