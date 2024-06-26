package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"

	"lystem/internal/agent"
	"lystem/internal/config"
	"lystem/internal/handlers"
	"lystem/internal/middleware"
	"lystem/pkg/postgres"
)

func main() {
	// ------- DATABASE -------
	db, err := postgres.NewStorage(config.Options.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	// ------- LOGGER -------
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLogger.Sync()

	// ------- CHANNELS TO SYNC GRACEFULLY SHUTDOWN -------
	appShutdownCh := make(chan os.Signal, 1)
	exitCh := make(chan os.Signal, 1)

	// ------- ORDERS INFO POLLER -------
	ctx, cancel := context.WithCancel(context.Background())
	ordersAgent := agent.New(db, config.Options, zapLogger)
	go ordersAgent.StartOrdersPolling(ctx)

	// ------- INIT APP -------
	app := fiber.New()

	// ------- HANDLERS -------
	v1 := handlers.New(db, ordersAgent)
	app.Use(logger.New(logger.Config{Output: os.Stdout}))
	api := app.Group("/api/user", middleware.Authorize(db))
	api.Post("/register", v1.CreateUser)
	api.Post("/login", v1.CreateSession)
	api.Delete("/logout", v1.DeleteSession)

	api.Post("/orders", v1.SaveOrder)
	api.Get("/orders", v1.GetOrders)

	api.Get("/balance", v1.GetBalance)
	api.Post("/balance/withdraw", v1.Withdraw)
	api.Get("/withdrawals", v1.Withdrawals)

	// ------- GRACEFULLY SHUTDOWN -------
	go func() {
		// wait until app.Listen got error and appShutdownCh got notification
		<-appShutdownCh

		// cancel context, that we sent to orders poller
		cancel()

		fmt.Println("Gracefully shutting down loystem application")
		// shutdown application
		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
		// gracefully close database connection
		db.Close()
		// notify exitCh to stop main goroutine
		signal.Notify(exitCh, os.Interrupt)

		// exit with code 0 to signal, that everything gone right
		os.Exit(0)
	}()

	// ------- START SERVER -------
	if err = app.Listen(config.Options.Address); err != nil {
		zapLogger.Error(err.Error())
		signal.Notify(appShutdownCh, os.Interrupt)
	}

	<-exitCh
}
