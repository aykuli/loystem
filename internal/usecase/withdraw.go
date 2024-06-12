package usecase

import (
	"errors"

	"github.com/valyala/fasthttp"

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

func (uc *WithdrawalUsecase) Create(ctx *fasthttp.RequestCtx, wRequest request.WithdrawRequest, currentUser *user.User) error {
	order, err := uc.db.FindOrderByNumber(ctx, wRequest.Order)
	if err != nil {
		return err
	}

	// does it needed to check order status
	if order.UserID != currentUser.ID {
		return ErrOrderUserIncorrect
	}

	userBalance, err := uc.db.FindBalance(ctx, currentUser)
	if err != nil {
		return err
	}
	if userBalance.Current < wRequest.Sum {
		return ErrNotEnoughBalance
	}

	return uc.db.CreateWithdraw(ctx, order, currentUser, wRequest.Sum)
}

func (uc *WithdrawalUsecase) FindAll(ctx *fasthttp.RequestCtx, currUser *user.User) ([]withdrawal.Withdrawal, error) {
	userBalance, err := uc.db.FindBalance(ctx, currUser)
	if err != nil {
		return nil, err
	}
	return uc.db.FindWithdrawals(ctx, userBalance)
}