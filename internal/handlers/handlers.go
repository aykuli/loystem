package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"lystem/internal/agent"
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

	//PollOrdersInfo()
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

// GetBalance godoc
//
//	@Summary		Получение текущего баланса пользователя
//	@Tags			Баланс
//	@Accept			text/plain
//	@Produce		application/json
//	@Success		200		{string}	json	"успешная обработка запроса"
//	@Failure		401		{string}	error	"пользователь не аутентифицирован"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/balance	[get]
func (v1 v1Handler) GetBalance(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(*user.User)
	usersUsecase := usecase.NewUserUsecase(v1.storage)
	balance, withdrawals, err := usersUsecase.GetBalance(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return ctx.JSON(presenter.NewBalanceResponse(balance, withdrawals))
}

// Withdraw godoc
//
//	@Summary		Получение текущего баланса пользователя
//	@Tags			Баланс
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200		{string}	json	"успешная обработка запроса"
//	@Failure		401		{string}	error	"пользователь не аутентифицирован"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/balance/withdraw	[post]
func (v1 v1Handler) Withdraw(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var wRequest request.WithdrawRequest
	if err := ctx.BodyParser(&wRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	if err := wRequest.Validate(); err != nil {
		fmt.Printf("\n\nvalidate error: %+v\n\n", err)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	currentUser := ctx.Locals("current_user").(*user.User)
	withdrawalUsecase := usecase.NewWithdrawalUsecase(v1.storage)

	err := withdrawalUsecase.Create(ctx.Context(), wRequest, currentUser)
	fmt.Printf("\n\n1 create error: %+v\n\n", err)

	if err != nil && errors.Is(err, usecase.ErrNotEnoughBalance) {
		return ctx.Status(fiber.StatusPaymentRequired).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil && errors.Is(err, usecase.ErrOrderUserIncorrect) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	return ctx.JSON(presenter.Common{Success: true})
}

// Withdrawals godoc
//
//	@Summary		Получение информации о выводе средств
//	@Tags			Списания
//	@Produce		application/json
//	@Success		200		{string}	json	"успешная обработка запроса"
//	@Success		204		{string}	json	"нет ни одного списания"
//	@Failure		401		{string}	error	"пользователь не аутентифицирован"
//	@Failure		500		{string}	error	"внутренняя ошибка сервера"
//	@Router			/api/user/withdrawals	    [get]
func (v1 v1Handler) Withdrawals(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("current_user").(*user.User)
	withdrawalsUsecase := usecase.NewWithdrawalUsecase(v1.storage)
	withdrawals, err := withdrawalsUsecase.FindAll(ctx.Context(), currentUser)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.Common{Success: false, Message: err.Error()})
	}

	if len(withdrawals) == 0 {
		return ctx.Status(fiber.StatusNoContent).JSON(presenter.Common{Success: true})
	}

	return ctx.JSON(presenter.NewWithdrawalsResponse(withdrawals))
}
