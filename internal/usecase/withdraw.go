package usecase

import (
	"context"
	"errors"

	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type WithdrawalUsecase struct {
	db storage.Storage
}

var (
	ErrNotEnoughBalance   = errors.New("not enough balance")
	ErrOrderUserIncorrect = errors.New("order user incorrect")
)

func NewWithdrawalUsecase(db storage.Storage) *WithdrawalUsecase {
	return &WithdrawalUsecase{db}
}

func (uc *WithdrawalUsecase) Create(ctx context.Context, wRequest request.WithdrawRequest, currentUser *user.User) (*withdrawal.Withdrawal, error) {
	userBalance, err := uc.db.FindBalance(ctx, currentUser)
	if err != nil {
		return nil, err
	}
	if userBalance.Current < wRequest.Sum {
		return nil, ErrNotEnoughBalance
	}

	return uc.db.CreateWithdraw(ctx, wRequest.Order, currentUser, wRequest.Sum)
}

func (uc *WithdrawalUsecase) FindAll(ctx context.Context, currUser *user.User) ([]withdrawal.Withdrawal, error) {
	userBalance, err := uc.db.FindBalance(ctx, currUser)
	if err != nil {
		return nil, err
	}
	return uc.db.FindWithdrawals(ctx, userBalance)
}
