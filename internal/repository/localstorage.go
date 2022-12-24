package repository

import "errors"

type ItemLocalStorageRepository struct {
	items map[string]string
}

func (s *ItemLocalStorageRepository) Add(id string, url string) error {
	s.items[id] = url
	return nil
}

func (s *ItemLocalStorageRepository) Get(id string) (string, error) {
	v, ok := s.items[id]
	if !ok {
		return "", errors.New("this url not found")
	}
	return v, nil
}

func NewItemLocalStorageRepository() *ItemLocalStorageRepository {
	return &ItemLocalStorageRepository{
		items: map[string]string{},
	}
}
