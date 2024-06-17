package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/loms/api/openapiv2"
	"route256/loms/internal/app/server"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/pkg/api/loms/v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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
		log.Printf("starting server app on url %s\n", config.GrpcUrl)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write(openapiv2.Doc)
	})
	mux.HandleFunc("/swaggerui/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
	mux.Handle("/", gwmux)

	gwServer := &http.Server{
		Addr:    config.HttpUrl,
		Handler: mux,
	}

	go func() {
		log.Printf("starting grpc-gateway on url %s", config.HttpUrl)
		log.Fatalln(gwServer.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("server shutting down")
	gwServer.Shutdown(context.Background())
	gprcServer.GracefulStop()
}
