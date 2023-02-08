package repositories

import "github.com/vdgalyns/link-shortener/internal/entities"

type Memory struct {
	links []entities.Link
}

func (m *Memory) Get(id string) (entities.Link, error) {
	for _, link := range m.links {
		if link.ID == id {
			return link, nil
		}
	}
	return entities.Link{}, ErrNotFound
}

func (m *Memory) GetByOriginalURL(originalURL string) (entities.Link, error) {
	for _, link := range m.links {
		if link.OriginalURL == originalURL {
			return link, nil
		}
	}
	return entities.Link{}, ErrNotFound
}

func (m *Memory) GetAllByUserID(userID string) ([]entities.Link, error) {
	links := make([]entities.Link, 0, len(m.links))
	for _, link := range m.links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}
	return links, nil
}

func (m *Memory) Add(link entities.Link) error {
	_, err := m.Get(link.ID)
	if err != nil {
		m.links = append(m.links, link)
	}
	return nil
}

func (m *Memory) Ping() error {
	return nil
}

func (m *Memory) AddBatch(links []entities.Link) error {
	for _, link := range links {
		err := m.Add(link)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMemory() *Memory {
	return &Memory{links: make([]entities.Link, 0)}
}
