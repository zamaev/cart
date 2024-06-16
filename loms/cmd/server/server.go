package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/loms/internal/app/server"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/pkg/api/loms/v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := config.NewConfig()

	lis, err := net.Listen("tcp", config.GrpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	gprcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Panic,
			middleware.Validate,
		),
	)
	reflection.Register(gprcServer)

	lomsServer := server.NewServer()
	loms.RegisterLomsServer(gprcServer, lomsServer)

	go func() {
		log.Println("starting server app")
		if err := gprcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	gwmux := runtime.NewServeMux()

	grpcClient, err := grpc.NewClient(config.GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	if err = loms.RegisterLomsHandler(context.Background(), gwmux, grpcClient); err != nil {
		log.Fatalln("failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    config.HttpUrl,
		Handler: gwmux,
	}

	go func() {
		log.Printf("starting grpc-gateway")
		log.Fatalln(gwServer.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("server shutting down")
	gwServer.Shutdown(context.Background())
	gprcServer.GracefulStop()
}
