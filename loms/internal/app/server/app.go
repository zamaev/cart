package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"route256/loms/api/openapiv2"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	GrpcServer *grpc.Server
	GwServer   *http.Server
	LomsServer *Server
}

func NewApp(ctx context.Context, config config.Config) *App {
	lomsServer := NewServer(config)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			middleware.Tracer,
			middleware.Metrics,
			middleware.Panic,
			middleware.Validate,
		),
	)
	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, lomsServer)

	grpcClient, err := grpc.NewClient(config.GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Panicw(ctx, "grpc.NewClient", "err", err)
	}
	gwmux := runtime.NewServeMux()
	if err = loms.RegisterLomsHandler(ctx, gwmux, grpcClient); err != nil {
		logger.Panicw(ctx, "loms.RegisterLomsHandler", "err", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) { w.Write(openapiv2.Doc) })
	mux.HandleFunc("/swaggerui/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))
	mux.Handle("/", gwmux)

	return &App{
		GrpcServer: grpcServer,
		GwServer: &http.Server{
			Addr:    config.HttpUrl,
			Handler: mux,
		},
		LomsServer: lomsServer,
	}
}

func (app *App) Shutdown(ctx context.Context) {
	logger.Infow(ctx, "shutting down loms server app")
	app.LomsServer.Shutdown()
}
