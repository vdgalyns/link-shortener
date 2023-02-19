package services

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"github.com/vdgalyns/link-shortener/internal/repositories"
)

type Links struct {
	repositories *repositories.Repositories
	config       *config.Config
}

func (l *Links) Get(hash string) (entities.Link, error) {
	_, err := entities.ValidateLinkHash(hash)
	if err != nil {
		if errors.Is(err, repositories.ErrLinkIsDeleted) {
			return entities.Link{}, ErrLinkIsDeleted
		}
		return entities.Link{}, err
	}
	return l.repositories.Get(hash)
}

func (l *Links) Add(originalURL, userID string) (string, error) {
	valid := entities.ValidateURL(originalURL)
	if !valid {
		return "", ErrLinkNotValid
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
	err = l.repositories.Add(link)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			existedLink, err := l.repositories.GetByOriginalURL(link.OriginalURL)
			if err != nil {
				return "", err
			}
			return l.config.BaseURL + "/" + existedLink.Hash, ErrLinkIsExist
		}
	}
	if err != nil {
		return "", err
	}
	readyURL := l.config.BaseURL + "/" + linkHash
	return readyURL, nil
}

func (l *Links) GetAllByUserID(userID string) ([]entities.Link, error) {
	_, err := entities.ValidateUserID(userID)
	if err != nil {
		return nil, err
	}
	return l.repositories.GetAllByUserID(userID)
}

func (l *Links) Ping() error {
	return l.repositories.Ping()
}

func (l *Links) AddBatch(originalURLs []string, userID string) ([]string, error) {
	links := make([]entities.Link, 0, len(originalURLs))
	for _, originalURL := range originalURLs {
		valid := entities.ValidateURL(originalURL)
		if !valid {
			return nil, ErrLinkNotValid
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
	if err := l.repositories.AddBatch(links); err != nil {
		return nil, err
	}
	output := make([]string, 0, len(originalURLs))
	for _, link := range links {
		output = append(output, l.config.BaseURL+"/"+link.Hash)
	}
	return output, nil
}

func (l *Links) RemoveBatch(urlHashes []string, userID string) error {
	return l.repositories.RemoveBatch(urlHashes, userID)
}

func NewLinks(repositories *repositories.Repositories, config *config.Config) *Links {
	return &Links{
		repositories: repositories,
		config:       config,
	}
}
