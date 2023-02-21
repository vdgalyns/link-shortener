package repositories

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/vdgalyns/link-shortener/internal/entities"
	"golang.org/x/sync/errgroup"
)

type Database struct {
	db *sql.DB
	mu *sync.RWMutex
}

func (d *Database) Get(hash string) (entities.Link, error) {
	if d.db == nil {
		return entities.Link{}, ErrDatabaseNotInitialized
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	var deletedAt sql.NullTime
	link := entities.Link{}
	row := d.db.QueryRow("SELECT hash, user_id, original_url, deleted_at FROM shortened_links WHERE hash = $1", hash)
	err := row.Scan(&link.Hash, &link.UserID, &link.OriginalURL, &deletedAt)
	if err != nil {
		return link, err
	}
	if deletedAt.Valid {
		return link, ErrLinkIsDeleted
	}
	return link, nil
}

func (d *Database) GetByOriginalURL(originalURL string) (entities.Link, error) {
	if d.db == nil {
		return entities.Link{}, ErrDatabaseNotInitialized
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
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
	d.mu.Lock()
	defer d.mu.Unlock()
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
	d.mu.RLock()
	defer d.mu.RUnlock()
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
	d.mu.Lock()
	defer d.mu.Unlock()
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

func (d *Database) RemoveBatch(urlHashes []string, userID string) error {
	if d.db == nil {
		return ErrDatabaseNotInitialized
	}
	d.mu.Lock()
	// defer d.mu.Unlock()
	g, ctx := errgroup.WithContext(context.Background())
	for _, urlHash := range urlHashes {
		urlHash := urlHash
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				d.mu.Lock()
				defer d.mu.Unlock()
				_, err := d.db.Exec(
					"UPDATE shortened_links SET deleted_at = $1 WHERE hash = $2 AND user_id = $3",
					time.Now(),
					urlHash,
					userID,
				)
				return err
			}
		})
	}

	go func() {
		g.Wait()
		d.mu.Unlock()
	}()

	return nil
}

func NewDatabase(database *sql.DB) *Database {
	return &Database{db: database, mu: &sync.RWMutex{}}
}
