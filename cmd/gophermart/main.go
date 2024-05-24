package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"lystem/internal/config"
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
}
