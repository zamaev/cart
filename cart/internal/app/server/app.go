package server

import (
	"context"
	"log"
	"net/http"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/middleware"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cart"
	"route256/cart/internal/pkg/service/loms"
	"route256/cart/internal/pkg/service/product"
	lomsapi "route256/cart/pkg/api/loms/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	http.Server
	Config     config.Config
	grpcClient *grpc.ClientConn
}

func NewApp() *App {
	config := config.NewConfig()

	grpcClient, err := grpc.NewClient(config.LomsServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	lomsClient := lomsapi.NewLomsClient(grpcClient)
	lomsService := loms.NewLomsService(lomsClient)

	cartRepository := repository.NewCartMemoryRepository()
	productService := product.NewProductService(config)
	cartService := cart.NewCartService(cartRepository, productService, lomsService)
	cartServer := NewServer(cartService)

	mux := http.NewServeMux()
	mux.Handle("POST /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.AddProduct))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.RemoveProduct))
	mux.Handle("DELETE /user/{user_id}/cart", middleware.ErrorWrapper(cartServer.ClearCart))
	mux.Handle("GET /user/{user_id}/cart/list", middleware.ErrorWrapper(cartServer.GetCart))
	mux.Handle("POST /cart/checkout", middleware.ErrorWrapper(cartServer.Checkout))

	h := middleware.LoggerWrapperHandler{
		Wrap: mux,
	}

	return &App{
		Server: http.Server{
			Addr:    config.CartServiceUrl,
			Handler: h,
		},
		Config:     config,
		grpcClient: grpcClient,
	}
}

func (app *App) ListenAndServe() error {
	log.Printf("starting server app on url %s\n", app.Config.CartServiceUrl)
	return app.Server.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	log.Println("shutting down server app")
	if err := app.grpcClient.Close(); err != nil {
		log.Printf("failed to close grpc client: %v\n", err)
	}
	return app.Server.Shutdown(ctx)
}
