package repositories

import (
	"context"
	"github.com/vdgalyns/link-shortener/internal/database"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Database struct {
	dsn string
}

func (d *Database) Get(hash string) (entities.URL, error) {
	return entities.URL{}, nil
}

func (d *Database) Add(url entities.URL) error {
	return nil
}

func (d *Database) GetAllByUserID(userID string) ([]entities.URL, error) {
	return nil, nil
}

func (d *Database) Ping() error {
	conn, err := database.NewDatabase(d.dsn)
	if err != nil {
		return err
	}
	return conn.Ping(context.Background())
}

func NewDatabase(dsn string) *Database {
	return &Database{dsn}
}
