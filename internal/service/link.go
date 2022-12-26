package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"strings"
)

var (
	ErrLinkIncorrect = errors.New("link incorrect")
)

type LinkService struct {
	repository repository.Link
}

func (s *LinkService) Get(hash string) (string, error) {
	valid := s.validateHash(hash)
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
	valid := s.validateURL(url)
	if !valid {
		return "", ErrLinkIncorrect
	}
	hash := s.createHash(url)
	err := s.repository.Add(hash, url)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *LinkService) validateURL(url string) bool {
	_, domain, _ := strings.Cut(url, "//")
	if len(domain) == 0 {
		return false
	}
	name, zone, _ := strings.Cut(domain, ".")
	if len(name) == 0 || len(zone) == 0 {
		return false
	}
	return true
}

func (s *LinkService) createHash(url string) string {
	data := []byte(url)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return hash
}

func (s *LinkService) validateHash(hash string) bool {
	data := []byte(hash)
	// TODO: умножение на 2, костыль (не знаю как по другому)
	return len(data) == md5.Size*2
}

func NewLinkService(repository repository.Link) *LinkService {
	return &LinkService{
		repository: repository,
	}
}
