package usecase

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/factory"
	"lystem/internal/models/session"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type CreateUseCase struct {
	db      storage.Storage
	factory *factory.UserFactory
}

func NewUserUseCase(db storage.Storage) *CreateUseCase {
	return &CreateUseCase{db, factory.NewUserFactory()}
}

func (uc *CreateUseCase) CreateUserAndSession(ctx *fasthttp.RequestCtx, userRequest request.CreateUser) (*session.Session, error) {
	newUser, err := uc.factory.Build(userRequest)
	if err != nil {
		return nil, err
	}

	savedUser, err := uc.db.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return uc.db.CreateSession(ctx, savedUser)
}
