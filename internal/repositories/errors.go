package repositories

import "errors"

var (
	ErrLinkNotFound           = errors.New("link not found")
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrLinkIsDeleted          = errors.New("link is deleted")
)
