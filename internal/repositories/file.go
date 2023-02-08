package repositories

import (
	"bufio"
	"encoding/json"
	"github.com/vdgalyns/link-shortener/internal/entities"
	"os"
)

type File struct {
	filePath string
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
	return url, ErrNotFound
}

func (f *File) readAllByPredicate(predicate string) ([]entities.Link, error) {
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
		if link.OriginalURL == predicate || link.Hash == predicate || link.UserID == predicate {
			links = append(links, link)
		}
	}
	return links, scanner.Err()
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
	return f.read(hash)
}

func (f *File) GetByOriginalURL(originalURL string) (entities.Link, error) {
	return f.read(originalURL)
}

func (f *File) GetAllByUserID(userID string) ([]entities.Link, error) {
	return f.readAllByPredicate(userID)
}

func (f *File) Add(link entities.Link) error {
	_, err := f.read(link.Hash)
	if err != nil {
		switch err {
		case ErrNotFound:
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

func NewFile(filePath string) *File {
	return &File{filePath}
}
