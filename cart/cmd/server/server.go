package main

import (
	"log"
	"net/http"
	"os"
	"route256/cart/internal/app/server"
	"route256/cart/internal/pkg/middleware"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cart"
	"route256/cart/internal/pkg/service/product"
)

func main() {
	cartRepository := repository.NewCartMemoryRepository()
	productService := product.NewProductService(os.Getenv("PRODUCT_SERVICE_TOKEN"))
	cartService := cart.NewCartService(cartRepository, productService)
	cartServer := server.NewServer(cartService)

	mux := http.NewServeMux()
	mux.Handle("POST /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.AddProduct))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", middleware.ErrorWrapper(cartServer.RemoveProduct))
	mux.Handle("DELETE /user/{user_id}/cart", middleware.ErrorWrapper(cartServer.ClearCart))
	mux.Handle("GET /user/{user_id}/cart/list", middleware.ErrorWrapper(cartServer.GetCart))

	h := middleware.LoggerWrapperHandler{
		Wrap: mux,
	}

	log.Println("starting server app")
	log.Fatal(http.ListenAndServe(":8082", h))
}
