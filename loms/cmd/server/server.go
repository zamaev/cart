package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"route256/loms/internal/app/server"
	"route256/loms/internal/pkg/config"
	"syscall"
)

func main() {
	config := config.NewConfig()
	app := server.NewApp(config)

	lis, err := net.Listen("tcp", config.GrpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("starting server app on url %s\n", config.GrpcUrl)
		if err := app.GrpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Printf("starting grpc-gateway on url %s", config.HttpUrl)
		if err := app.GwServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("server shutting down")
	if err := app.GwServer.Shutdown(context.Background()); err != nil {
		log.Printf("grpc-gateway shutdown error: %v\n", err)
	}
	app.GrpcServer.GracefulStop()
}
