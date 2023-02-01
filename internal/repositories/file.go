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

func (f *File) read(id string) (entities.URL, error) {
	url := entities.URL{}
	file, err := f.open()
	if err != nil {
		return url, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		foundUrl := entities.URL{}
		if err = json.Unmarshal(scanner.Bytes(), &foundUrl); err != nil {
			return url, err
		}
		if foundUrl.ID == id {
			return foundUrl, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return url, err
	}
	return url, ErrNotFound
}

func (f *File) readAllByPredicate(predicate string) ([]entities.URL, error) {
	urls := make([]entities.URL, 0)
	file, err := f.open()
	if err != nil {
		return urls, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := entities.URL{}
		if err = json.Unmarshal(scanner.Bytes(), &url); err != nil {
			return urls, err
		}
		if url.OriginalURL == predicate || url.ID == predicate || url.UserID == predicate {
			urls = append(urls, url)
		}
	}
	err = scanner.Err()
	return urls, err
}

func (f *File) write(url entities.URL) error {
	file, err := f.open()
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(url)
}

func (f *File) Get(id string) (entities.URL, error) {
	return f.read(id)
}

func (f *File) GetAllByUserId(userId string) ([]entities.URL, error) {
	return f.readAllByPredicate(userId)
}

func (f *File) Add(url entities.URL) error {
	_, err := f.read(url.ID)
	if err != nil {
		switch err {
		case ErrNotFound:
			return f.write(url)
		default:
			return err
		}
	}
	return nil
}

func NewFile(filePath string) *File {
	return &File{filePath}
}