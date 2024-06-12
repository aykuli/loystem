package handlers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"lystem/internal/agent"
	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/presenter"
	"lystem/internal/request"
	"lystem/internal/usecase"
)

// SaveOrder godoc
//
//	@Summary		Загрузка заказа пользователем
//	@Tags			Заказ
//	@Accept			text/plain
//	@Produce		application/json
//	@Param			payload	body		string
//	@Success		200		{string}	json	"номер заказа уже был загружен этим пользователем"
//	@Success		202		{string}	json	"новый номер заказа принят в обработку"
//	@Failure		400		{string}	error	"неверный формат запроса"
//	@Failure		401		{string}	error	"пользователь не аутентифицирован"
//	@Failure		409		{string}	error	"номер заказы был уже загружен другим пользователем"
//	@Failure		422		{string}	error	"неверный формат номера заказа"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/orders	[post]
func (v1 v1Handler) SaveOrder(ctx *fiber.Ctx) error {
	var saveOrderRequest request.SaveOrderRequest
	if err := saveOrderRequest.Parse(ctx.Body()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}
	if err := saveOrderRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	currentUser := ctx.Locals("current_user").(*user.User)
	orderUsecase := usecase.NewOrderUsecase(v1.storage)

	//check if order already registered in system.
	foundOrder, err := orderUsecase.FindByNumber(ctx.Context(), saveOrderRequest.Number)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if foundOrder != nil {
		// if it registered with the same user - 200
		if foundOrder.UserID == currentUser.ID {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "order already registered for this user"})
		}
		// if it registered with the other user - 409
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"success": true, "message": "order already registered by other user"})
	}

	newOrder, err := orderUsecase.Save(ctx.Context(), saveOrderRequest, currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{})
	}

	go v1.getOrderInfo(ctx.Context(), newOrder, currentUser)

	return ctx.Status(fiber.StatusAccepted).JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) getOrderInfo(ctx *fasthttp.RequestCtx, newOrder *order.Order, currentUser *user.User) {
	orderUsecase := usecase.NewOrderUsecase(v1.storage)

	orderWithInfo, err := v1.agent.GetOrderInfo(newOrder)
	if orderWithInfo == nil {
		log.Printf("[GET ORDER INFO] order with number is empty %s\n err: %s\n", newOrder.Number, err)
		return
	}
	if err != nil && errors.Is(err, agent.ErrNoSuchOrder) {
		newOrder.Status = order.StatusInvalid
		newOrder.Accrual = 0
		err = orderUsecase.Update(ctx, newOrder, currentUser)
		if err != nil {
			log.Printf("[Usecase] update failed for order with number %s\n err: %s\n", newOrder.Number, err)
		}
		return
	} else if err != nil {
		log.Printf("[Agent] request failed for order with number %s\n err: %s\n", newOrder.Number, err)
		return
	}
	err = orderUsecase.Update(ctx, orderWithInfo, currentUser)
	if err != nil {
		log.Printf("[Usecase] update failed for order with number %s\n err: %s\n", newOrder.Number, err)
	}
}

// GetOrders godoc
//
//	@Summary		Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
//	@Tags			Заказ
//	@Accept			text/plain
//	@Produce		application/json
//	@Success		200		{string}	json	"успешная обработка запроса"
//	@Success		204		{string}	json	"нет данных для ответа"
//	@Failure		401		{string}	error	"пользователь не аутентифицирован"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/orders	[get]
func (v1 v1Handler) GetOrders(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(*user.User)
	orderUsecase := usecase.NewOrderUsecase(v1.storage)
	orders, err := orderUsecase.FindAll(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	if len(orders) == 0 {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"success": true, "message": "no orders found"})
	}

	return ctx.JSON(presenter.NewOrdersResponse(orders))
}

//func (v1 v1Handler) PollOrdersInfo() {
//	orderUsecase := usecase.NewOrderUsecase(v1.storage)
//	orders, err := orderUsecase.FindAllProcessable(v1.storage)
//
//	ordersTicker := time.NewTicker(config.Options.PollInterval)
//	defer ordersTicker.Stop()
//	for {
//		select {
//		case <-ordersTicker.C:
//			ordersTicker.Reset()
//			//v1.
//		}
//	}
//
//}
