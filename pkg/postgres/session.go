package postgres

import (
	"context"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/repository"
)

func (s *DBStorage) CreateSession(ctx context.Context, currentUser *user.User) (*session.Session, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	sessionsRepo := repository.NewSessionsRepository(conn)

	// transaction used for 2 db table changing operations - creating and deleting sessions
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	newSession, err := sessionsRepo.Create(ctx, tx, currentUser)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return newSession, nil
}

func (s *DBStorage) DeleteSession(ctx context.Context, u *user.User) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	sessionsRepo := repository.NewSessionsRepository(conn)

	err = sessionsRepo.Delete(ctx, u)
	if err != nil {
		return newDBError(err)
	}

	return nil
}
