package middleware

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"

	"lystem/internal/presenter"
	"lystem/pkg/postgres"
)

var ignorePaths = []string{
	"/api/user/login",
	"/api/user/register",
}

var errUnauththorized = errors.New("не удалось идентифицировать пользователя")

func Authorize(ctx *fiber.Ctx) error {

	if ctx.Method() == fiber.MethodPost && slices.Contains(ignorePaths, ctx.Path()) {
		return ctx.Next()
	}

	var token string
	authHeader := ctx.Get(fiber.HeaderAuthorization)
	if authHeader != "" {
		token = extractToken(authHeader)
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Payload: errUnauththorized.Error()})
		}
	}

	u, err := postgres.FindUserByToken(ctx.Context(), token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.Common{Success: false, Payload: errUnauththorized.Error()})
	}

	type currentUser string
	type currentToken string
	ctx.SetUserContext(context.WithValue(ctx.Context(), currentUser("current_user"), *u))
	ctx.SetUserContext(context.WithValue(ctx.UserContext(), currentToken("current_token"), token))

	return ctx.Next()
}

// Extracts token from header value
// Example: "Authorization": "Token token=f0bd6f98-4771-11ee-be56-0242ac120002"
func extractToken(headerValue string) string {
	parts := strings.Split(headerValue, "Token token=")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}
