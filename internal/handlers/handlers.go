package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

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

	GetOrders(ctx *fiber.Ctx) error
	GetBalance(ctx *fiber.Ctx) error
	Withdraw(ctx *fiber.Ctx) error
	Withdrawals(ctx *fiber.Ctx) error
}

func New(db storage.Storage) Handler {
	return v1Handler{storage: db}
}

type v1Handler struct {
	storage storage.Storage
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
	fmt.Println("-------------\nPOST api/user/register CreateUser\n----------------")
	var userRequest request.CreateUser
	if err := ctx.BodyParser(&userRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}
	if err := userRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}

	userUseCase := usecase.NewUserUseCase(v1.storage)
	session, err := userUseCase.CreateUserAndSession(ctx.Context(), userRequest)
	if err != nil && errors.Is(err, postgres.ErrUserAlreadyExists) {
		return ctx.Status(fiber.StatusConflict).JSON(presenter.Common{Success: false, Payload: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}

	bearerToken := fmt.Sprintf("Token token=%s", session.ID)
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
//	@Router			/api/user/login	[post]
func (v1 v1Handler) CreateSession(ctx *fiber.Ctx) error {
	fmt.Println("-------------\nPOST api/user/login CreateSession\n----------------")

	var sessionRequest request.CreateSession
	if err := ctx.BodyParser(&sessionRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}
	if err := sessionRequest.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}

	sessionUseCase := usecase.NewSessionUseCase(v1.storage)
	session, err := sessionUseCase.Create(ctx.Context(), sessionRequest)
	if err != nil && errors.Is(err, usecase.ErrInvalidCreds) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Payload: err.Error()})
	} else if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Payload: err.Error()})
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
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Payload: "invalid creds"})
	}

	sessionUseCase := usecase.NewSessionUseCase(v1.storage)
	if err := sessionUseCase.Delete(ctx.Context(), currentUser); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.Common{Success: false, Payload: err.Error()})
	}

	ctx.Context().RemoveUserValue("current_token")
	ctx.Context().RemoveUserValue("current_user")
	return ctx.JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) GetBalance(ctx *fiber.Ctx) error {
	return ctx.JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) Withdraw(ctx *fiber.Ctx) error {
	return ctx.JSON(presenter.Common{Success: true})
}

func (v1 v1Handler) Withdrawals(ctx *fiber.Ctx) error {
	return ctx.JSON(presenter.Common{Success: true})
}
