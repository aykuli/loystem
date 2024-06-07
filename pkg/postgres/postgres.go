package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/repository"
)

var (
	Instance *pgxpool.Pool
	pgOnce   sync.Once
)

func FindUserByToken(ctx context.Context, token string) (*user.User, error) {
	conn, err := Instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	sessionsRepo := repository.NewSessionsRepository(conn)
	usersRepo := repository.NewUsersRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	session, err := sessionsRepo.FindByID(ctx, tx, token)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	foundUser, err := usersRepo.FindByID(session.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, rollbackOnErr(ctx, tx, err)
	} else if err == nil && foundUser != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, newDBError(rollbackErr)
		}
		return nil, ErrUserAlreadyExists
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}

	return foundUser, nil
}

type DBStorage struct {
}

func NewStorage(uri string) (*DBStorage, error) {
	fmt.Println("--------\nNewStorage\n-----------")

	ctx := context.Background()
	var s *DBStorage
	if err := s.createDBPool(ctx, uri); err != nil {
		return nil, err
	}
	if err := s.init(ctx); err != nil {
		return &DBStorage{}, err
	}

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

		Instance = pool
	})

	if resErr != nil {
		return resErr
	}

	return nil
}

func (s *DBStorage) init(ctx context.Context) error {
	conn, err := Instance.Acquire(ctx)
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
		if err = tx.Rollback(ctx); err != nil {
			return newDBError(err)
		}
		return newDBError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}

	return nil

}

func (s *DBStorage) CreateUser(ctx *fasthttp.RequestCtx, user *user.User) (*user.User, error) {
	conn, err := Instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	usersRepo := repository.NewUsersRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	// find out if such user already exists
	foundUser, err := usersRepo.FindByLogin(user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, rollbackOnErr(ctx, tx, err)
	} else if err == nil && foundUser != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, newDBError(rollbackErr)
		}
		return nil, ErrUserAlreadyExists
	}

	savedUser, err := usersRepo.Create(ctx, tx, user)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return savedUser, nil
}

func (s *DBStorage) FindUserByLogin(ctx *fasthttp.RequestCtx, login string) (*user.User, error) {
	conn, err := Instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	usersRepo := repository.NewUsersRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	foundUser, err := usersRepo.FindByLogin(login)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, newDBError(rollbackErr)
		}
		return nil, pgx.ErrNoRows
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}

	return foundUser, nil
}

func (s *DBStorage) CreateSession(ctx *fasthttp.RequestCtx, user *user.User) (*session.Session, error) {
	conn, err := Instance.Acquire(ctx)
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

	newSession, err := sessionsRepo.Create(ctx, tx, user)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}

	return newSession, nil
}

func (s *DBStorage) DeleteSession(ctx *fasthttp.RequestCtx, u *user.User) error {
	conn, err := Instance.Acquire(ctx)
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

func rollbackOnErr(ctx context.Context, tx pgx.Tx, err error) error {
	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		return newDBError(rollbackErr)
	}

	return newDBError(err)
}
