package service

import (
	"github.com/vdgalyns/link-shortener/internal/repository"
)

type Item interface {
	Add(url string) (string, error)
	Get(id string) (string, error)
}

type Service struct {
	Item
}

func NewService(repositories *repository.Repository) *Service {
	return &Service{
		Item: NewItemService(repositories.Item),
	}
}
