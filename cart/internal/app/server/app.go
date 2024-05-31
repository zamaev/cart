package server

import (
	"log"
	"net/http"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/middleware"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cart"
	"route256/cart/internal/pkg/service/product"
)

type App struct {
	http.Server
	Config config.Config
}

func NewApp() *App {
	config := config.NewConfig()
	cartRepository := repository.NewCartMemoryRepository()
	productService := product.NewProductService(config)
	cartService := cart.NewCartService(cartRepository, productService)
	cartServer := NewServer(cartService)

	mux := http.NewServeMux()
	mux.Handle("POST /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.AddProduct))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.RemoveProduct))
	mux.Handle("DELETE /user/{user_id}/cart", middleware.ErrorWrapper(cartServer.ClearCart))
	mux.Handle("GET /user/{user_id}/cart/list", middleware.ErrorWrapper(cartServer.GetCart))

	h := middleware.LoggerWrapperHandler{
		Wrap: mux,
	}

	return &App{
		Server: http.Server{
			Addr:    ":8082",
			Handler: h,
		},
		Config: config,
	}
}

func (app *App) ListenAndServe() error {
	log.Println("starting server app")
	return app.Server.ListenAndServe()
}
