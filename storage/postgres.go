package storage

import (
	"context"
	"database/sql"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(dataSourseName string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dataSourseName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func (s *PostgresStorage) Close() error {
	return s.DB.Close()
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}
