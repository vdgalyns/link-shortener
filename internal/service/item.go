package service

import (
	"crypto/md5"
	"fmt"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"io"
)

type ItemService struct {
	repository repository.Item
}

func (s *ItemService) Add(url string) (string, error) {
	h := md5.New()
	io.WriteString(h, url)
	sum := h.Sum(nil)
	id := fmt.Sprintf("%x", sum)
	err := s.repository.Add(id, url)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *ItemService) Get(id string) (string, error) {
	v, err := s.repository.Get(id)
	return v, err
}

func NewItemService(repository repository.Item) *ItemService {
	return &ItemService{
		repository: repository,
	}
}
