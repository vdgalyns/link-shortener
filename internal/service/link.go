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

const hashLen int = 6

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
	if len(name) == 0 || len(zone) < 2 {
		return false
	}
	return true
}

func (s *LinkService) createHash(url string) string {
	data := []byte(url)
	sum := fmt.Sprintf("%x", md5.Sum(data))
	hash := strings.Builder{}
	for i, v := range sum {
		if i < hashLen {
			hash.WriteString(string(v))
			continue
		}
		break
	}
	return hash.String()
}

func (s *LinkService) validateHash(hash string) bool {
	return len(hash) == hashLen
}

func NewLinkService(repository repository.Link) *LinkService {
	return &LinkService{
		repository: repository,
	}
}
