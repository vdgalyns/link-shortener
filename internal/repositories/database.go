package repositories

import (
	"database/sql"
	"fmt"

	"github.com/vdgalyns/link-shortener/internal/entities"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Get(hash string) (entities.URL, error) {
	if d.db == nil {
		return entities.URL{}, ErrDatabaseNotInitialized
	}
	url := entities.URL{}
	row := d.db.QueryRow("SELECT hash, user_id, original_url FROM urls WHERE hash = $1", hash)
	err := row.Scan(&url.Hash, &url.UserID, &url.OriginalURL)
	if err != nil {
		fmt.Println(err)
		return entities.URL{}, err
	}
	return url, nil
}

func (d *Database) Add(url entities.URL) error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	_, err := d.db.Exec(
		"INSERT INTO urls (hash, user_id, original_url) VALUES($1, $2, $3)",
		url.Hash,
		url.UserID,
		url.OriginalURL,
	)
	return err
}

func (d *Database) GetAllByUserID(userID string) ([]entities.URL, error) {
	if d.db == nil {
		return nil, ErrDatabaseNotInitialized
	}
	rows, err := d.db.Query("SELECT hash, user_id, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := make([]entities.URL, 0)
	for rows.Next() {
		var url entities.URL
		err = rows.Scan(&url.Hash, &url.UserID, &url.OriginalURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (d *Database) Ping() error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	return d.db.Ping()
}

func (d *Database) AddBatch(urls []entities.URL) error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO urls(hash,user_id,original_url) VALUES($1,$2,$3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, v := range urls {
		if _, err = stmt.Exec(v.Hash, v.UserID, v.OriginalURL); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func NewDatabase(database *sql.DB) *Database {
	return &Database{db: database}
}
