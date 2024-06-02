package usecase

import (
	"github.com/valyala/fasthttp"

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

func (uc *OrderUsecase) FindByNumber(ctx *fasthttp.RequestCtx, number string) (*order.Order, error) {
	return uc.db.FindOrderByNumber(ctx, number)
}

func (uc *OrderUsecase) Save(ctx *fasthttp.RequestCtx, req request.SaveOrderRequest, currUser *user.User) (*order.Order, error) {
	return uc.db.SaveOrder(ctx, req.Number, currUser.ID)
}

func (uc *OrderUsecase) Update(ctx *fasthttp.RequestCtx, newOrder *order.Order, currUser *user.User) error {
	savedOrder, err := uc.db.UpdateOrder(ctx, newOrder)
	if err != nil {
		return err
	}

	if savedOrder.Status != order.StatusProcessed {
		return nil
	}

	err = uc.db.IncreaseBalance(ctx, savedOrder, currUser)
	if err != nil {
		savedOrder.Status = order.StatusProcessing
		_, err = uc.db.UpdateOrder(ctx, savedOrder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *OrderUsecase) FindAll(ctx *fasthttp.RequestCtx, currentUser *user.User) ([]order.Order, error) {
	return uc.db.SelectOrders(ctx, currentUser)

}
