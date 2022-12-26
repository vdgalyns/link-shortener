package repository

import "errors"

var (
	ErrLinkNotFound = errors.New("link not found")
)

type LinkLocalStorageRepository struct {
	links map[string]string
}

func (r *LinkLocalStorageRepository) Get(hash string) (string, error) {
	url, ok := r.links[hash]
	if !ok {
		return "", ErrLinkNotFound
	}
	return url, nil
}

func (r *LinkLocalStorageRepository) Add(hash string, url string) error {
	_, err := r.Get(hash)
	if err != nil {
		r.links[hash] = url
		return nil
	}
	return nil
}

func NewLinkLocalStorageRepository() *LinkLocalStorageRepository {
	return &LinkLocalStorageRepository{
		links: map[string]string{},
	}
}
