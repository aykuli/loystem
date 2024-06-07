package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/repository"
)

type DBStorage struct {
	instance *pgxpool.Pool
}

func NewStorage(uri string) (*DBStorage, error) {
	ctx := context.Background()
	var s *DBStorage
	tryCount := 0
	createConn := func() error {
		word := "try"
		if tryCount > 0 {
			word = "retry"
		}
		fmt.Printf("%s to connect to database, probe %d\n", word, tryCount)
		tryCount++

		pool, err := pgxpool.New(ctx, uri)
		if err != nil {
			return fmt.Errorf("could not connect to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			return fmt.Errorf("could not connect to database: %v", err)
		}

		fmt.Printf("  | -- connected to database %s\n", uri)
		s.instance = pool

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 12 * time.Second
	if err := backoff.Retry(createConn, expBackoff); err != nil {
		return nil, fmt.Errorf("\nfailed to connect to database after retrying %d times: %v", tryCount, err)
	}
	if err := s.init(ctx); err != nil {
		return &DBStorage{}, err
	}

	return s, nil
}

func (s *DBStorage) init(ctx context.Context) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	repo := repository.New(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	err = repo.Init(ctx, tx)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}

	return nil
}

func (s *DBStorage) Close() {
	s.instance.Close()
}
