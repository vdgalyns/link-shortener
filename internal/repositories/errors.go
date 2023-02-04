package repositories

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)
