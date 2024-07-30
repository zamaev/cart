package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"route256/cart/internal/pkg/cache"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/middleware"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cart"
	"route256/cart/internal/pkg/service/loms"
	"route256/cart/internal/pkg/service/product"
	"route256/cart/internal/pkg/service/product/product_cache"
	lomsapi "route256/cart/pkg/api/loms/v1"
	"route256/cart/pkg/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	http.Server
	config     config.Config
	grpcClient *grpc.ClientConn
}

func NewApp(ctx context.Context, config config.Config) *App {
	grpcClient, err := grpc.NewClient(
		config.LomsServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Panicw(ctx, "grpc.NewClient", "err", err)
	}
	lomsClient := lomsapi.NewLomsClient(grpcClient)
	lomsService := loms.NewLomsService(lomsClient)

	cartRepository := repository.NewCartMemoryRepository()
	productService := product.NewProductService(config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Panicw(ctx, "cache.Ping", "err", err)
	}
	cache := cache.NewRedisLRUCache(redisClient, config.CacheSize)
	productCacheService := product_cache.NewProductCacheService(productService, cache, config.CacheDefaultTTL)

	cartService := cart.NewCartService(cartRepository, productCacheService, lomsService)
	cartServer := NewServer(cartService)

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	muxTracerWrapper := middleware.NewMuxTracerWrapper(mux)
	muxMetricsWrapper := middleware.NewMuxMetricsWrapper(muxTracerWrapper)
	muxMetricsWrapper.Handle("POST /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.AddProduct))
	muxMetricsWrapper.Handle("DELETE /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.RemoveProduct))
	muxMetricsWrapper.Handle("DELETE /user/{user_id}/cart", middleware.ErrorWrapper(cartServer.ClearCart))
	muxMetricsWrapper.Handle("GET /user/{user_id}/cart/list", middleware.ErrorWrapper(cartServer.GetCart))
	muxMetricsWrapper.Handle("POST /cart/checkout", middleware.ErrorWrapper(cartServer.Checkout))

	return &App{
		Server: http.Server{
			Addr:    config.CartServiceUrl,
			Handler: mux,
		},
		config:     config,
		grpcClient: grpcClient,
	}
}

func (app *App) ListenAndServe(ctx context.Context) error {
	logger.Infow(ctx, "starting server app", "url", app.config.CartServiceUrl)
	return app.Server.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	logger.Infow(ctx, "shutting down server app")
	if err := app.grpcClient.Close(); err != nil {
		logger.Errorw(ctx, "failed to close grpc client", "err", err)
	}
	return app.Server.Shutdown(ctx)
}
