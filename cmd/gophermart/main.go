package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// ------- CONTEXT & WAIT GROUP FOR SYNC AGENT GOROUTINE GRACEFULLY SHUTDOWN WITH APPLICATION'S & DB -------
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// ------- ORDERS INFO POLLER -------
	ordersAgent := agent.New(db, config.Options, zapLogger)
	wg.Add(1)
	go ordersAgent.StartOrdersPolling(ctx, &wg)

	// ------- INIT APP -------
	app := fiber.New()

	// ------- HANDLERS -------
	v1 := handlers.New(db, ordersAgent, config.Options)
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
	exit := make(chan os.Signal, 1)
	waiter := make(chan os.Signal, 1)
	go func() {
		// waiter wait until it gets one of:
		//   |-- application interrupting signals on err | syscall.SIGTERM -- terminated
		//   |-- keyboard Ctrl+C                         | syscall.SIGINT  -- interrupt
		signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT) // terminated

		// blocks here until there's a signal
		sig := <-waiter
		zapLogger.Info("1 Signal notify received: " + sig.String())

		// cancel context, that we sent to orders poller
		cancel()

		zapLogger.Info("2 Gracefully shutting down loystem application")

		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
		zapLogger.Info("3 Application shut down")

		wg.Wait()
		zapLogger.Info("5 Finish waiting agent wait group done()")

		// gracefully close database connection after finishing agent work
		db.Close()
		zapLogger.Info("6 Database connections closed")

		// signal main goroutine that gracefully shutdown finished
		exit <- syscall.SIGSTOP
	}()

	// ------- START SERVER -------
	if err = app.Listen(config.Options.Address); err != nil {
		zapLogger.Error(err.Error())
		waiter <- syscall.SIGTERM
	}

	<-exit
	zapLogger.Info("7 Main goroutine exited.")
	// exit with code 0 to signal, that everything gone right
	os.Exit(0)
}
