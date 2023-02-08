package services

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"github.com/vdgalyns/link-shortener/internal/repositories"
)

type Urls struct {
	repositories *repositories.Repositories
	config       *config.Config
}

func (u *Urls) Get(hash string) (entities.Link, error) {
	_, err := entities.ValidateLinkHash(hash)
	if err != nil {
		return entities.Link{}, err
	}
	return u.repositories.Get(hash)
}

func (u *Urls) Add(originalURL, userID string) (string, error) {
	valid := entities.ValidateURL(originalURL)
	if !valid {
		return "", ErrURLNotValid
	}
	linkHash, err := entities.CreateLinkHash(originalURL)
	if err != nil {
		return "", err
	}
	link := entities.Link{
		Hash:        linkHash,
		UserID:      userID,
		OriginalURL: originalURL,
	}
	err = u.repositories.Add(link)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			existedLink, err := u.repositories.GetByOriginalURL(link.OriginalURL)
			if err != nil {
				return "", err
			}
			return u.config.BaseURL + "/" + existedLink.Hash, ErrURLIsExist
		}
	}
	if err != nil {
		return "", err
	}
	readyURL := u.config.BaseURL + "/" + linkHash
	return readyURL, nil
}

func (u *Urls) GetAllByUserID(userID string) ([]entities.Link, error) {
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
	links := make([]entities.Link, 0, len(originalURLs))
	for _, originalURL := range originalURLs {
		valid := entities.ValidateURL(originalURL)
		if !valid {
			return nil, ErrURLNotValid
		}
		linkHash, err := entities.CreateLinkHash(originalURL)
		if err != nil {
			return nil, err
		}
		link := entities.Link{
			Hash:        linkHash,
			UserID:      userID,
			OriginalURL: originalURL,
		}
		links = append(links, link)
	}
	if err := u.repositories.AddBatch(links); err != nil {
		return nil, err
	}
	output := make([]string, 0, len(originalURLs))
	for _, link := range links {
		output = append(output, u.config.BaseURL+"/"+link.Hash)
	}
	return output, nil
}

func NewUrls(repositories *repositories.Repositories, config *config.Config) *Urls {
	return &Urls{
		repositories: repositories,
		config:       config,
	}
}