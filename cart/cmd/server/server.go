package main

import (
	"context"
	"os"
	"os/signal"
	"route256/cart/internal/app/server"
	"route256/cart/internal/pkg/config"
	"route256/cart/pkg/logger"
	"route256/cart/pkg/tracing"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	ctx := context.Background()
	config := config.NewConfig()

	logger.Set(logger.With("service", config.ServiceName))

	tp, err := tracing.NewTracerProvider(ctx, config.ServiceName, config.TracerUrl)
	if err != nil {
		logger.Panicw(ctx, "tracing.NewTracerProvider", "err", err)
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracing.Set(tp.Tracer(config.ServiceName))
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Errorw(ctx, "tp.Shutdown", "err", err)
		}
	}()

	app := server.NewApp(ctx, config)
	go func() {
		if err := app.ListenAndServe(ctx); err != nil {
			logger.Panicw(ctx, "app.ListenAndServe", "err", err)
		}
	}()
	defer func() {
		logger.Infow(ctx, "app shutting down")
		if err := app.Shutdown(ctx); err != nil {
			logger.Errorw(ctx, "app.Shutdown", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
