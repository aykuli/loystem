package middleware

import (
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"

	"lystem/internal/presenter"
	"lystem/internal/storage"
)

var ignorePaths = []string{
	"/api/user/login",
	"/api/user/register",
}

var errUnauththorized = errors.New("не удалось идентифицировать пользователя")

func Authorize(db storage.Storage) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		if ctx.Method() == fiber.MethodPost && slices.Contains(ignorePaths, ctx.Path()) {
			return ctx.Next()
		}

		var token string
		authHeader := ctx.Get(fiber.HeaderAuthorization)
		if authHeader != "" {
			token = extractToken(authHeader)
			if token == "" {
				return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.NewFailure(errUnauththorized))
			}
		}

		foundUser, err := db.FindUserByToken(ctx.Context(), token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.NewFailure(errUnauththorized))
		}

		ctx.Locals("current_user", foundUser)
		ctx.Locals("current_token", token)
		return ctx.Next()
	}
}

// Extracts token from header value
// Example: "Authorization": "Token token=<session-id-as-token>"
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
