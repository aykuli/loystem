package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"lystem/internal/agent"
	"lystem/internal/config"
	"lystem/internal/handlers"
	"lystem/internal/middleware"
	"lystem/pkg/postgres"
)

func main() {
	fmt.Printf("config: %+v\n\n", config.Options)
	db, err := postgres.NewStorage(config.Options.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down loystem application")
		db.Close()

		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	app.Use(logger.New(logger.Config{Output: os.Stdout}))

	v1 := handlers.New(db, agent.New())

	api := app.Group("/api/user", middleware.Authorize(db))
	api.Post("/register", v1.CreateUser)
	api.Post("/login", v1.CreateSession)
	api.Delete("/logout", v1.DeleteSession)

	api.Post("/orders", v1.SaveOrder)
	api.Get("/orders", v1.GetOrders)

	api.Get("/balance", v1.GetBalance)
	api.Post("/balance/withdraw", v1.Withdraw)
	api.Post("/withdrawals", v1.Withdrawals)

	if err := app.Listen(config.Options.Address); err != nil {
		c <- os.Interrupt
	}
}
