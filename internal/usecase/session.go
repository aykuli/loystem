package usecase

import (
	"errors"

	"github.com/valyala/fasthttp"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type SessionUseCase struct {
	db storage.Storage
}

var ErrInvalidCreds = errors.New("неверная пара логин/пароль")

func NewSessionUseCase(db storage.Storage) *SessionUseCase {
	return &SessionUseCase{db}
}

func (uc *SessionUseCase) Create(ctx *fasthttp.RequestCtx, sessionRequest request.CreateSession) (*session.Session, error) {
	user, err := uc.db.FindUserByLogin(ctx, sessionRequest.Login)
	if err != nil {
		return nil, err
	}

	valid := user.ValidatePassword(sessionRequest.Password)
	if !valid {
		return nil, ErrInvalidCreds
	}

	newSession, err := uc.db.CreateSession(ctx, user)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (uc *SessionUseCase) Delete(ctx *fasthttp.RequestCtx, u *user.User) error {
	return uc.db.DeleteSession(ctx, u)
}
