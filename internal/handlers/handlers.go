package handlers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"lystem/internal/agent"
	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/presenter"
	"lystem/internal/request"
	"lystem/internal/storage"
	"lystem/internal/usecase"
	"lystem/pkg/postgres"
)

type Handler interface {
	CreateUser(ctx *fiber.Ctx) error
	CreateSession(ctx *fiber.Ctx) error
	DeleteSession(ctx *fiber.Ctx) error

	SaveOrder(ctx *fiber.Ctx) error
	GetOrders(ctx *fiber.Ctx) error
	GetBalance(ctx *fiber.Ctx) error
	Withdraw(ctx *fiber.Ctx) error
	Withdrawals(ctx *fiber.Ctx) error
}

func New(db storage.Storage, agent *agent.GophermartAgent) Handler {
	return v1Handler{storage: db, agent: agent}
}

type v1Handler struct {
	storage storage.Storage
	agent   *agent.GophermartAgent
}

// CreateUser godoc
//
//	@Summary		Создание пользователя
//	@Tags			Пользователи
//	@Accept			application/json
//	@Produce		application/json
//	@Param			payload	body		request.CreateUser
//	@Success		200		{string}	json	"пользователь успешно зарегестрирован и аутентифицирован"
//	@Failure		400		{string}	error	"неверный формат запроса"
//	@Failure		409		{string}	error	"логин уже занят"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/register	[post]
func (v1 v1Handler) CreateUser(ctx *fiber.Ctx) error {
	var userRequest request.CreateUser
	if err := ctx.BodyParser(&userRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}
	if err := userRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	userUsecase := usecase.NewUserUsecase(v1.storage)
	newSession, err := userUsecase.CreateUserAndSession(ctx.Context(), userRequest)
	if err != nil && errors.Is(err, postgres.ErrUserAlreadyExists) {
		return ctx.Status(fiber.StatusConflict).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	bearerToken := fmt.Sprintf("Token token=%s", newSession.ID)
	ctx.Set("Authorization", bearerToken)

	return ctx.JSON(presenter.Common{Success: true})
}

// CreateSession godoc
//
//	@Summary		Аутентификация пользователя
//	@Tags			Сессия
//	@Accept			application/json
//	@Produce		application/json
//	@Param			payload	body		request.CreateSession
//	@Success		200		{string}	json	"пользователь успешно аутентифицирован"
//	@Failure		400		{string}	error	"неверный формат запроса"
//	@Failure		401		{string}	error	"неверная пара логин/пароль"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/login	    [post]
func (v1 v1Handler) CreateSession(ctx *fiber.Ctx) error {
	var sessionRequest request.CreateSession
	if err := ctx.BodyParser(&sessionRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}
	if err := sessionRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	sessionUsecase := usecase.NewSessionUsecase(v1.storage)
	session, err := sessionUsecase.Create(ctx.Context(), sessionRequest)
	if err != nil && errors.Is(err, usecase.ErrInvalidCreds) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	bearerToken := fmt.Sprintf("Token token=%s", session.ID)
	ctx.Set("Authorization", bearerToken)

	return ctx.JSON(presenter.Common{Success: true})
}

// DeleteSession godoc
//
//	@Summary		Удаление сессии
//	@Tags			Авторизация
//	@Accept			application/json
//	@Produce		application/json
//	@Param			payload	body		request.CreateSession
//	@Success		200		{string}	json	"сессия успешно удалена"
//	@Failure		401		{string}	error	"не удалось идентифицировать пользователя"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/logout	[post]
func (v1 v1Handler) DeleteSession(ctx *fiber.Ctx) error {
	currentUser, ok := ctx.UserContext().Value("current_user").(*user.User)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Message: "invalid creds"})
	}

	sessionUsecase := usecase.NewSessionUsecase(v1.storage)
	if err := sessionUsecase.Delete(ctx.Context(), currentUser); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	ctx.Context().RemoveUserValue("current_token")
	ctx.Context().RemoveUserValue("current_user")
	return ctx.JSON(presenter.Common{Success: true})
}

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
//	@Router			/api/user/orders	[post]
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

func (v1 v1Handler) GetBalance(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(*user.User)
	usersUsecase := usecase.NewUserUsecase(v1.storage)
	balance, withdrawals, err := usersUsecase.GetBalance(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return ctx.JSON(presenter.NewBalanceResponse(balance, withdrawals))
}

func (v1 v1Handler) Withdraw(ctx *fiber.Ctx) error {
	var wRequest request.WithdrawRequest
	if err := ctx.BodyParser(&wRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}
	if err := wRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	currentUser := ctx.Locals("current_user").(*user.User)
	withdrawalUsecase := usecase.NewWithdrawalUsecase(v1.storage)

	err := withdrawalUsecase.Create(ctx.Context(), wRequest, currentUser)
	if err != nil && errors.Is(err, usecase.ErrNotEnoughBalance) {
		return ctx.Status(fiber.StatusPaymentRequired).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil && errors.Is(err, usecase.ErrOrderUserIncorrect) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	return ctx.JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) Withdrawals(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(*user.User)
	withdrawalsUsecase := usecase.NewWithdrawalUsecase(v1.storage)
	withdrawals, err := withdrawalsUsecase.FindAll(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	if len(withdrawals) == 0 {
		return ctx.Status(fiber.StatusAccepted).JSON(presenter.Common{Success: true})
	}
	return ctx.JSON(presenter.NewWithdrawalsResponse(withdrawals))
}
