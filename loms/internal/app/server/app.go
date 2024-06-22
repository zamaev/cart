package server

import (
	"context"
	"log"
	"net/http"
	"route256/loms/api/openapiv2"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/pkg/api/loms/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	GrpcServer *grpc.Server
	GwServer   *http.Server
}

func NewApp(config config.Config) *App {
	lomsServer := NewServer(config)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Panic,
			middleware.Validate,
		),
	)
	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, lomsServer)

	grpcClient, err := grpc.NewClient(config.GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	gwmux := runtime.NewServeMux()
	if err = loms.RegisterLomsHandler(context.Background(), gwmux, grpcClient); err != nil {
		log.Fatalln("failed to register gateway:", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) { w.Write(openapiv2.Doc) })
	mux.HandleFunc("/swaggerui/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))
	mux.Handle("/", gwmux)

	return &App{
		GrpcServer: grpcServer,
		GwServer: &http.Server{
			Addr:    config.HttpUrl,
			Handler: mux,
		},
	}
}
