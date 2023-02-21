package repositories

import (
	"github.com/vdgalyns/link-shortener/internal/entities"
	"sync"
	"time"
)

type Memory struct {
	links []entities.Link
	mu    *sync.RWMutex
}

func (m *Memory) Get(hash string) (entities.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, link := range m.links {
		if link.Hash == hash {
			if link.DeletedAt != "" {
				return link, ErrLinkIsDeleted
			}
			return link, nil
		}
	}
	return entities.Link{}, ErrLinkNotFound
}

func (m *Memory) GetByOriginalURL(originalURL string) (entities.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, link := range m.links {
		if link.OriginalURL == originalURL {
			return link, nil
		}
	}
	return entities.Link{}, ErrLinkNotFound
}

func (m *Memory) GetAllByUserID(userID string) ([]entities.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	links := make([]entities.Link, 0, len(m.links))
	for _, link := range m.links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}
	return links, nil
}

func (m *Memory) Add(link entities.Link) error {
	_, err := m.Get(link.Hash)
	m.mu.Lock()
	defer m.mu.Unlock()
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

func (m *Memory) RemoveBatch(urlHashes []string, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, urlHash := range urlHashes {
		for i := range m.links {
			if m.links[i].Hash == urlHash && m.links[i].UserID == userID {
				m.links[i].DeletedAt = time.Now().String()
				break
			}
		}
	}
	return nil
}

func NewMemory() *Memory {
	return &Memory{links: make([]entities.Link, 0), mu: &sync.RWMutex{}}
}
