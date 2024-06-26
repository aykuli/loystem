package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"lystem/internal/models/user"
	"lystem/internal/repository"
)

func (s *DBStorage) CreateUser(ctx context.Context, newUser *user.User) (*user.User, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	usersRepo := repository.NewUsersRepository(conn)
	balancesRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	// find out if such user already exists
	foundUser, err := usersRepo.FindByLogin(newUser.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, rollbackOnErr(ctx, tx, err)
	} else if err == nil && foundUser != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, newDBError(rollbackErr)
		}
		return nil, ErrUserAlreadyExists
	}

	savedUser, err := usersRepo.Create(ctx, tx, newUser)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if _, err = balancesRepo.Create(ctx, tx, savedUser); err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return savedUser, nil
}

func (s *DBStorage) FindUserByLogin(ctx context.Context, login string) (*user.User, error) {
	conn, err := s.instance.Acquire(ctx)
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

func (s *DBStorage) FindUserByToken(ctx context.Context, token string) (*user.User, error) {
	conn, err := s.instance.Acquire(ctx)
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

	foundSession, err := sessionsRepo.FindByID(ctx, tx, token)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	} else if foundSession == nil {
		return nil, rollbackOnErr(ctx, tx, errors.New("session not found"))
	}

	foundUser, err := usersRepo.FindByID(foundSession.UserID)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	} else if foundUser == nil {
		return nil, rollbackOnErr(ctx, tx, errors.New("user not found"))
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}

	return foundUser, nil
}
