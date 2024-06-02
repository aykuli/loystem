package usecase

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/factory"
	"lystem/internal/models/balance"
	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type UserUsecase struct {
	db      storage.Storage
	factory *factory.UserFactory
}

func NewUserUsecase(db storage.Storage) *UserUsecase {
	return &UserUsecase{db, factory.NewUserFactory()}
}

func (uc *UserUsecase) CreateUserAndSession(ctx *fasthttp.RequestCtx, req request.CreateUser) (*session.Session, error) {
	newUser, err := uc.factory.Build(req)
	if err != nil {
		return nil, err
	}

	savedUser, err := uc.db.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return uc.db.CreateSession(ctx, savedUser)
}

func (uc *UserUsecase) GetBalance(ctx *fasthttp.RequestCtx, currUser *user.User) (*balance.Balance, []withdrawal.Withdrawal, error) {
	var withdrawals []withdrawal.Withdrawal
	userBalance, err := uc.db.FindBalance(ctx, currUser)
	if err != nil {
		return nil, withdrawals, err
	}

	withdrawals, err = uc.db.FindWithdrawals(ctx, userBalance)
	if err != nil {
		return userBalance, withdrawals, err
	}

	return userBalance, withdrawals, nil
}
