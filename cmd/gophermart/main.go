package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"lystem/internal/config"
	"lystem/internal/handlers"
	"lystem/internal/middleware"
	"lystem/pkg/postgres"
)

func main() {
	_, err := postgres.NewStorage(config.Options.DatabaseUri)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))
	err = app.Listen(config.Options.Address)
	if err != nil {
		log.Fatal(err)
	}
	v1 := handlers.New()
	api := app.Group("/api/user", middleware.Authorize)
	// /api/v1/register
	api.Post("/register", v1.CreateUser)
	// /api/v1/login
	api.Post("/login", v1.CreateSession)
	// /api/v1/logout
	api.Delete("/logout", v1.DeleteSession)
}
