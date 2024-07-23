package main

import (
	"context"
	"os"
	"os/signal"
	"route256/notifier/internal/app"
	"route256/notifier/internal/pkg/config"
	"route256/notifier/pkg/logger"
	"sync"
	"syscall"
)

func main() {
	ctx, close := context.WithCancel(context.Background())

	config := config.NewConfig()

	app, err := app.NewApp(config)
	if err != nil {
		logger.Panicw(ctx, "app.NewApp", "err", err)
	}

	wg := &sync.WaitGroup{}
	app.Run(ctx, wg)
	defer app.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	close()
	wg.Wait()
}
