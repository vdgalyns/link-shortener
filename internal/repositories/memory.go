package repositories

import "github.com/vdgalyns/link-shortener/internal/entities"

type Memory struct {
	urls []entities.URL
}

func (m *Memory) Get(hash string) (entities.URL, error) {
	for _, url := range m.urls {
		if url.Hash == hash {
			return url, nil
		}
	}
	return entities.URL{}, ErrNotFound
}

func (m *Memory) GetAllByUserId(userId string) ([]entities.URL, error) {
	urls := make([]entities.URL, 0)
	for _, url := range m.urls {
		if url.UserID == userId {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

func (m *Memory) Add(url entities.URL) error {
	_, err := m.Get(url.Hash)
	if err != nil {
		m.urls = append(m.urls, url)
	}
	return nil
}

func NewMemory() *Memory {
	return &Memory{urls: make([]entities.URL, 0)}
}
