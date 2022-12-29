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
	repository repository.Link
}

func (s *LinkService) Get(hash string) (string, error) {
	valid := utils.ValidateHash(hash)
	if !valid {
		return "", ErrLinkIncorrect
	}
	url, err := s.repository.Get(hash)
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
	err := s.repository.Add(hash, url)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func NewLinkService(repository repository.Link) *LinkService {
	return &LinkService{
		repository: repository,
	}
}
