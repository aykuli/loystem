package usecase

import (
	"errors"
	"fmt"

	"github.com/valyala/fasthttp"

	"lystem/internal/models/order"
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
	foundOrder, err := uc.db.FindOrderByNumber(ctx, wRequest.Order)
	fmt.Printf("\n\nfoundOrder: %+v\n", foundOrder)
	fmt.Printf("err: %+v\n", err)
	if err != nil {
		return err
	}
	if foundOrder == nil || foundOrder.Status != order.StatusProcessed {
		return ErrOrderUserIncorrect
	}
	// does it needed to check order status
	if foundOrder.UserID != currentUser.ID {
		return ErrOrderUserIncorrect
	}

	userBalance, err := uc.db.FindBalance(ctx, currentUser)
	if err != nil {
		return err
	}
	if userBalance.Current < wRequest.Sum {
		return ErrNotEnoughBalance
	}

	return uc.db.CreateWithdraw(ctx, foundOrder, currentUser, wRequest.Sum)
}

func (uc *WithdrawalUsecase) FindAll(ctx *fasthttp.RequestCtx, currUser *user.User) ([]withdrawal.Withdrawal, error) {
	userBalance, err := uc.db.FindBalance(ctx, currUser)
	if err != nil {
		return nil, err
	}
	return uc.db.FindWithdrawals(ctx, userBalance)
}
