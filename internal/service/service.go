package service

import (
	"github.com/vdgalyns/link-shortener/internal/repository"
)

type Link interface {
	Get(hash string) (string, error)
	Add(url string) (string, error)
}

type Service struct {
	Link
}

func NewService(repositories *repository.Repository) *Service {
	return &Service{
		Link: NewLinkService(repositories),
	}
}
