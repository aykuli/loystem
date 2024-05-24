package usecase

import (
	"context"

	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/request"
	"lystem/internal/storage"
)

type OrderUsecase struct {
	db storage.Storage
}

func NewOrderUsecase(db storage.Storage) *OrderUsecase {
	return &OrderUsecase{db}
}

func (uc *OrderUsecase) FindByNumber(ctx context.Context, number string) (*order.Order, error) {
	return uc.db.FindOrderByNumber(ctx, number)
}

func (uc *OrderUsecase) Save(ctx context.Context, req request.SaveOrderRequest, currUser *user.User) (*order.Order, error) {
	return uc.db.SaveOrder(ctx, req.Number, currUser.ID)
}

func (uc *OrderUsecase) Update(ctx context.Context, newOrder *order.Order) error {
	savedOrder, err := uc.db.UpdateOrder(ctx, newOrder)
	if err != nil {
		return err
	}

	if savedOrder.Status != order.StatusProcessed {
		return nil
	}

	if err = uc.db.IncreaseBalance(ctx, savedOrder); err != nil {
		savedOrder.Status = order.StatusProcessing
		if _, err = uc.db.UpdateOrder(ctx, savedOrder); err != nil {
			return err
		}
	}

	return nil
}

func (uc *OrderUsecase) FindAllAccrual(ctx context.Context) ([]order.Order, error) {
	return uc.db.SelectAccrualOrders(ctx)
}

func (uc *OrderUsecase) FindUserOrders(ctx context.Context, u *user.User) ([]order.Order, error) {
	return uc.db.SelectUserOrders(ctx, u)
}
