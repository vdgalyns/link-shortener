package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDatabase(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
