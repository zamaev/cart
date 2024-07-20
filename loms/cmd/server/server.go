package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"route256/loms/internal/app/server"
	"route256/loms/internal/pkg/config"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/tracing"
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

	lis, err := net.Listen("tcp", config.GrpcUrl)
	if err != nil {
		logger.Panicw(ctx, "net.Listen", "err", err)
	}

	go func() {
		logger.Infow(ctx, "starting server app", "url", config.GrpcUrl)
		if err := app.GrpcServer.Serve(lis); err != nil {
			logger.Panicw(ctx, "app.GrpcServer.Serve", "err", err)
		}
	}()
	defer func() {
		logger.Infow(ctx, "shutting down grpc server app")
		app.GrpcServer.GracefulStop()
	}()

	go func() {
		logger.Infow(ctx, "starting grpc-gateway", "url", config.HttpUrl)
		if err := app.GwServer.ListenAndServe(); err != nil {
			logger.Panicw(ctx, "app.GwServer.ListenAndServe", "err", err)
		}
	}()
	defer func() {
		logger.Infow(ctx, "shutting down http server app")
		if err := app.GwServer.Shutdown(context.Background()); err != nil {
			logger.Errorw(ctx, "app.GwServer.Shutdown", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
