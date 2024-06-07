package usecase

import (
	"errors"

	"github.com/valyala/fasthttp"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type SessionUsecase struct {
	db storage.Storage
}

var ErrInvalidCreds = errors.New("неверная пара логин/пароль")

func NewSessionUsecase(db storage.Storage) *SessionUsecase {
	return &SessionUsecase{db}
}

func (uc *SessionUsecase) Create(ctx *fasthttp.RequestCtx, sessionRequest request.CreateSession) (*session.Session, error) {
	foundUser, err := uc.db.FindUserByLogin(ctx, sessionRequest.Login)
	if err != nil {
		return nil, err
	}

	valid := foundUser.ValidatePassword(sessionRequest.Password)

	if !valid {
		return nil, ErrInvalidCreds
	}

	newSession, err := uc.db.CreateSession(ctx, foundUser)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (uc *SessionUsecase) Delete(ctx *fasthttp.RequestCtx, u *user.User) error {
	return uc.db.DeleteSession(ctx, u)
}
