package usecase

import (
	"context"
	"errors"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type SessionUsecase struct {
	db   storage.Storage
	salt string
}

var ErrInvalidCreds = errors.New("неверная пара логин/пароль")

func NewSessionUsecase(db storage.Storage, salt string) *SessionUsecase {
	return &SessionUsecase{db, salt}
}

func (uc *SessionUsecase) Create(ctx context.Context, sessionRequest request.CreateSession) (*session.Session, error) {
	foundUser, err := uc.db.FindUserByLogin(ctx, sessionRequest.Login)
	if err != nil {
		return nil, err
	}

	valid := foundUser.ValidatePassword(sessionRequest.Password, uc.salt)
	if !valid {
		return nil, ErrInvalidCreds
	}

	newSession, err := uc.db.CreateSession(ctx, foundUser)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (uc *SessionUsecase) Delete(ctx context.Context, u *user.User) error {
	return uc.db.DeleteSession(ctx, u)
}
