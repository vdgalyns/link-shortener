package repositories

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrNotDatabase = errors.New("not database")
)
