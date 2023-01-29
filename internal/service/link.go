package service

import (
	"errors"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/utils"
)

var (
	ErrLinkIncorrect = errors.New("link incorrect")
)

type LinkService struct {
	repositories *repository.Repository
}

func (s *LinkService) Get(hash string) (string, error) {
	valid := utils.ValidateHash(hash)
	if !valid {
		return "", ErrLinkIncorrect
	}
	url, err := s.repositories.Get(hash)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *LinkService) Add(url string) (string, error) {
	valid := utils.ValidateURL(url)
	if !valid {
		return "", ErrLinkIncorrect
	}
	hash := utils.CreateHash(url)
	err := s.repositories.Add(hash, url)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func NewLinkService(repositories *repository.Repository) *LinkService {
	return &LinkService{
		repositories: repositories,
	}
}
