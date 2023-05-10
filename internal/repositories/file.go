package repositories

import (
	"bufio"
	"encoding/json"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"os"
	"sync"
	"time"
)

type File struct {
	filePath string
	mu       *sync.RWMutex
}

func (f *File) open() (*os.File, error) {
	return os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}

func (f *File) read(hash string) (entities.Link, error) {
	url := entities.Link{}
	file, err := f.open()
	if err != nil {
		return url, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		foundLink := entities.Link{}
		if err = json.Unmarshal(scanner.Bytes(), &foundLink); err != nil {
			return url, err
		}
		if foundLink.Hash == hash {
			return foundLink, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return url, err
	}
	return url, ErrLinkNotFound
}

func (f *File) readAll() ([]entities.Link, error) {
	links := make([]entities.Link, 0)
	file, err := f.open()
	if err != nil {
		return links, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := entities.Link{}
		if err = json.Unmarshal(scanner.Bytes(), &link); err != nil {
			return links, err
		}
		links = append(links, link)
	}
	return links, scanner.Err()
}

func (f *File) readAllByPredicate(predicate string) ([]entities.Link, error) {
	links, err := f.readAll()
	if err != nil {
		return []entities.Link{}, err
	}
	suitableLinks := make([]entities.Link, 0)
	for _, link := range links {
		if link.OriginalURL == predicate || link.Hash == predicate || link.UserID == predicate {
			suitableLinks = append(suitableLinks, link)
		}
	}
	return suitableLinks, nil
}

func (f *File) write(link entities.Link) error {
	file, err := f.open()
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(link)
}

func (f *File) Get(hash string) (entities.Link, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	link, err := f.read(hash)
	if err != nil {
		return entities.Link{}, err
	}
	if link.DeletedAt != "" {
		return entities.Link{}, ErrLinkIsDeleted
	}
	return link, nil
}

func (f *File) GetByOriginalURL(originalURL string) (entities.Link, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.read(originalURL)
}

func (f *File) GetAllByUserID(userID string) ([]entities.Link, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.readAllByPredicate(userID)
}

func (f *File) Add(link entities.Link) error {
	f.mu.RLock()
	_, err := f.read(link.Hash)
	f.mu.RUnlock()
	if err != nil {
		switch err {
		case ErrLinkNotFound:
			f.mu.Lock()
			defer f.mu.Unlock()
			return f.write(link)
		default:
			return err
		}
	}
	return nil
}

func (f *File) Ping() error {
	return nil
}

func (f *File) AddBatch(links []entities.Link) error {
	for _, link := range links {
		err := f.Add(link)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) RemoveBatch(urlHashes []string, userID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	file, err := os.OpenFile(f.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	links, err := f.readAll()
	if err != nil {
		return err
	}
	for i := range links {
		for _, urlHash := range urlHashes {
			if links[i].Hash == urlHash && links[i].UserID == userID {
				links[i].DeletedAt = time.Now().String()
				break
			}
		}
	}
	err = os.Truncate(f.filePath, 0)
	if err != nil {
		return err
	}
	for _, link := range links {
		err = f.write(link)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewFile(filePath string) *File {
	return &File{filePath: filePath, mu: &sync.RWMutex{}}
}
