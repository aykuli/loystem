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
	if newOrder.Status == order.StatusProcessed {
		err := uc.db.UpdateOrderAndIncreaseBalance(ctx, newOrder)
		return err
	}

	err := uc.db.UpdateOrder(ctx, newOrder)
	return err
}

func (uc *OrderUsecase) SelectUnprocessed(ctx context.Context, limit int) ([]order.Order, error) {
	return uc.db.SelectUnprocessedOrders(ctx, limit)
}

func (uc *OrderUsecase) FindAllUserOrders(ctx context.Context, u *user.User) ([]order.Order, error) {
	return uc.db.FindAllUserOrders(ctx, u)
}
