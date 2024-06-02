package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

var ignorePaths = []string{
	"/api/v1/login",
	"/api/v1/sign_in",
	"/api/v1/register",
}

func Authorize(ctx *fiber.Ctx) error {
	if ctx.Method() == fiber.MethodPost && slices.Contains(ignorePaths, ctx.Path()) {
		return ctx.Next()
	}

	return nil
}
