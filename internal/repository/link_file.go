package repository

import (
	"bufio"
	"encoding/json"
	"os"
)

type LinkFileRepository struct {
	fileName string
}

type LinkInFile struct {
	Hash, URL string
}

func (r *LinkFileRepository) open() (*os.File, error) {
	return os.OpenFile(r.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}

func (r *LinkFileRepository) read(hash string) (string, error) {
	file, err := r.open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var l LinkInFile
		err = json.Unmarshal(scanner.Bytes(), &l)
		if err != nil {
			return "", err
		}
		if l.Hash == hash {
			return l.URL, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", ErrLinkNotFound
}

func (r *LinkFileRepository) write(hash string, url string) error {
	file, err := r.open()
	if err != nil {
		return err
	}
	defer file.Close()
	l := LinkInFile{
		Hash: hash,
		URL:  url,
	}
	return json.NewEncoder(file).Encode(&l)
}

func (r *LinkFileRepository) Get(hash string) (string, error) {
	return r.read(hash)
}

func (r *LinkFileRepository) Add(hash string, url string) error {
	_, err := r.read(hash)
	if err == ErrLinkNotFound {
		return r.write(hash, url)
	}
	return nil
}

func NewLinkFileRepository(fileName string) *LinkFileRepository {
	return &LinkFileRepository{fileName}
}
