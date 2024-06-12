package postgres

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/repository"
)

func (s *DBStorage) CreateSession(ctx *fasthttp.RequestCtx, currentUser *user.User) (*session.Session, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	sessionsRepo := repository.NewSessionsRepository(conn)

	//find user
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

func (s *DBStorage) DeleteSession(ctx *fasthttp.RequestCtx, u *user.User) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	sessionsRepo := repository.NewSessionsRepository(conn)

	//find user
	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	err = sessionsRepo.Delete(ctx, tx, u)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}

	return nil
}
