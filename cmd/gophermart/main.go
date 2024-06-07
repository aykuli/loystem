package main

import (
	"fmt"
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
	db, err := postgres.NewStorage(config.Options.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))

	api := app.Group("/api/user", middleware.AcquireDBConnection, middleware.Authorize)

	v1 := handlers.New(db)
	api.Post("/register", v1.CreateUser)
	api.Post("/login", v1.CreateSession)
	api.Delete("/logout", v1.DeleteSession)

	api.Get("/orders", v1.GetOrders)
	api.Get("/balance", v1.GetBalance)
	api.Post("/balance/withdraw", v1.Withdraw)
	api.Post("/withdrawals", v1.Withdrawals)
	fmt.Println(config.Options)
	err = app.Listen(config.Options.Address)
	if err != nil {
		log.Fatal(err)
	}
}
