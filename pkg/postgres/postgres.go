package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	instance *pgxpool.Pool
	pgOnce   sync.Once
)

type DBStorage struct {
	instance *pgxpool.Pool
}

func NewStorage(uri string) (*DBStorage, error) {
	ctx := context.Background()
	var s *DBStorage
	if err := s.createDBPool(ctx, uri); err != nil {
		return &DBStorage{}, err
	}
	//if err := s.createMetricsTable(ctx); err != nil {
	//	return &DBStorage{}, err
	//}

	return s, nil
}

func (s *DBStorage) createDBPool(ctx context.Context, uri string) error {
	var resErr error
	pgOnce.Do(func() {
		pool, err := pgxpool.New(ctx, uri)
		if err != nil {
			resErr = err
			return
		}

		instance = pool
	})

	if resErr != nil {
		return resErr
	}

	return nil
}
