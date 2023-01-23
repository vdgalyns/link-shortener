package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type LinkFileRepository struct {
	fileName string
}

type FileItem struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

var (
	ErrFileNotRead = errors.New("file not read")
)

func (r *LinkFileRepository) openFile() (*os.File, error) {
	return os.OpenFile(r.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}

func (r *LinkFileRepository) readLink(hash string) (string, error) {
	file, err := r.openFile()
	if err != nil {
		return "", ErrFileNotRead
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var item FileItem
		if err = json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return "", err
		}
		if item.Hash == hash {
			return item.URL, err
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", ErrLinkNotFound
}

func (r *LinkFileRepository) writeLink(hash string, url string) error {
	file, err := r.openFile()
	if err != nil {
		return ErrFileNotRead
	}
	defer file.Close()
	item := FileItem{
		Hash: hash,
		URL:  url,
	}
	return json.NewEncoder(file).Encode(&item)
}

func (r *LinkFileRepository) Add(hash string, url string) error {
	_, err := r.Get(hash)
	if err == ErrLinkNotFound {
		return r.writeLink(hash, url)
	}
	return err
}

func (r *LinkFileRepository) Get(hash string) (string, error) {
	url, err := r.readLink(hash)
	if err != nil {
		return "", err
	}
	return url, nil
}

func NewLinkFileRepository(fileName string) *LinkFileRepository {
	return &LinkFileRepository{
		fileName,
	}
}
