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
	u.repositories.Get(hash)
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
	u.repositories.Add(url)
	if err != nil {
		return "", err
	}
	readyURL := u.config.BaseURL + "/" + hash
	return readyURL, nil
}

func (u *Urls) GetAllByUserID(userID string) ([]entities.URL, error) {
	_, err := entities.ValidateUserID(userID)
	if err != nil {
		return nil, err
	}
	return u.repositories.GetAllByUserID(userID)
}

func (u *Urls) Ping() error {
	return u.repositories.Ping()
}

func (u *Urls) AddBatch(originalURLs []string, userID string) ([]string, error) {
	urls := make([]entities.URL, 0, len(originalURLs))
	for _, v := range originalURLs {
		valid := entities.ValidateURL(v)
		if !valid {
			return nil, ErrURLNotValid
		}
		hash, err := entities.CreateURLHash(v)
		if err != nil {
			return nil, err
		}
		urls = append(urls, entities.URL{OriginalURL: v, UserID: userID, Hash: hash})
	}
	err := u.repositories.AddBatch(urls)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0, len(originalURLs))
	for _, v := range urls {
		output = append(output, u.config.BaseURL+"/"+v.Hash)
	}
	return output, nil
}

func NewUrls(repositories *repositories.Repositories, config *config.Config) *Urls {
	return &Urls{
		repositories: repositories,
		config:       config,
	}
}
