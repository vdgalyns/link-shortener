package repository

import "errors"

var (
	ErrLinkNotFound = errors.New("link not found")
)

type LinkMemoryRepository struct {
	links map[string]string
}

func (r *LinkMemoryRepository) Get(hash string) (string, error) {
	url, ok := r.links[hash]
	if !ok {
		return "", ErrLinkNotFound
	}
	return url, nil
}

func (r *LinkMemoryRepository) Add(hash string, url string) error {
	_, err := r.Get(hash)
	if err != nil {
		r.links[hash] = url
		return nil
	}
	return nil
}

func NewLinkMemoryRepository() *LinkMemoryRepository {
	return &LinkMemoryRepository{
		links: map[string]string{},
	}
}
