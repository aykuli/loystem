package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	CreateUser(ctx *fiber.Ctx) error
	CreateSession(ctx *fiber.Ctx) error
	DeleteSession(ctx *fiber.Ctx) error
}

func New() Handler {
	return V1Handler{}
}

type V1Handler struct {
}

func (h V1Handler) CreateUser(ctx *fiber.Ctx) error {
	return nil
}

func (h V1Handler) CreateSession(ctx *fiber.Ctx) error {
	return nil
}

func (h V1Handler) DeleteSession(ctx *fiber.Ctx) error {
	return nil
}
