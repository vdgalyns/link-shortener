package repositories

import (
	"database/sql"

	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Get(hash string) (entities.Link, error) {
	if d.db == nil {
		return entities.Link{}, ErrDatabaseNotInitialized
	}
	link := entities.Link{}
	row := d.db.QueryRow("SELECT hash, user_id, original_url FROM shortened_links WHERE hash = $1", hash)
	err := row.Scan(&link.Hash, &link.UserID, &link.OriginalURL)
	if err != nil {
		return entities.Link{}, err
	}
	return link, nil
}

func (d *Database) GetByOriginalURL(originalURL string) (entities.Link, error) {
	if d.db == nil {
		return entities.Link{}, ErrDatabaseNotInitialized
	}
	link := entities.Link{}
	row := d.db.QueryRow("SELECT hash, user_id, original_url FROM shortened_links WHERE original_url = $1", originalURL)
	err := row.Scan(&link.Hash, &link.UserID, &link.OriginalURL)
	if err != nil {
		return entities.Link{}, err
	}
	return link, nil
}

func (d *Database) Add(link entities.Link) error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	_, err := d.db.Exec(
		`INSERT INTO shortened_links (hash, user_id, original_url) VALUES($1, $2, $3)`,
		link.Hash,
		link.UserID,
		link.OriginalURL,
	)
	return err
}

func (d *Database) GetAllByUserID(userID string) ([]entities.Link, error) {
	if d.db == nil {
		return nil, ErrDatabaseNotInitialized
	}
	rows, err := d.db.Query("SELECT hash, user_id, original_url FROM shortened_links WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	links := make([]entities.Link, 0)
	for rows.Next() {
		var link entities.Link
		if err = rows.Scan(&link.Hash, &link.UserID, &link.OriginalURL); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (d *Database) Ping() error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	return d.db.Ping()
}

func (d *Database) AddBatch(links []entities.Link) error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO shortened_links (hash, user_id, original_url) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, link := range links {
		if _, err = stmt.Exec(link.Hash, link.UserID, link.OriginalURL); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func NewDatabase(database *sql.DB) *Database {
	return &Database{db: database}
}
