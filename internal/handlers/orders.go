package handlers

import (
	"github.com/gofiber/fiber/v2"

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
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.NewFailure(err))
	}
	if err := saveOrderRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.NewFailure(err))
	}

	currentUser := ctx.Locals("current_user").(*user.User)
	orderUsecase := usecase.NewOrderUsecase(v1.storage)

	//check if order already registered in system.
	foundOrder, err := orderUsecase.FindByNumber(ctx.Context(), saveOrderRequest.Number)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.NewFailure(err))
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.NewFailure(err))
	}

	return ctx.Status(fiber.StatusAccepted).JSON(presenter.NewSuccess(newOrder))
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
	orders, err := orderUsecase.FindAllUserOrders(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.NewFailure(err))
	}

	if len(orders) == 0 {
		return ctx.Status(fiber.StatusNoContent).JSON(presenter.NewSuccess([]order.Order{}))
	}

	ctx.Set("Content-Type", "application/json")
	return ctx.Status(fiber.StatusOK).JSON(presenter.NewOrdersResponse(orders))
}
