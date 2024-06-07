package middleware

import (
	"github.com/gofiber/fiber/v2"

	"lystem/pkg/postgres"
)

func AcquireDBConnection(ctx *fiber.Ctx) error {
	context := ctx.Context()
	conn, err := postgres.Instance.Acquire(context)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	context.SetUserValue("dbConnection", conn)
	defer func() {
		context.RemoveUserValue("dbConnection")
		conn.Release()
	}()

	return ctx.Next()
}
