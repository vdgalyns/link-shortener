package repository

import "errors"

var (
	ErrLinkNotFound = errors.New("link not found")
)

type LinkLocalRepository struct {
	links map[string]string
}

func (r *LinkLocalRepository) Get(hash string) (string, error) {
	url, ok := r.links[hash]
	if !ok {
		return "", ErrLinkNotFound
	}
	return url, nil
}

func (r *LinkLocalRepository) Add(hash string, url string) error {
	_, err := r.Get(hash)
	if err != nil {
		r.links[hash] = url
		return nil
	}
	return nil
}

func NewLinkLocalRepository() *LinkLocalRepository {
	return &LinkLocalRepository{
		links: map[string]string{},
	}
}
