package config

import "os"

type Config struct {
	CartServiceUrl      string
	LomsServiceUrl      string
	ProductServiceToken string
}

func NewConfig() Config {
	cartServiceUrl := os.Getenv("CART_SERVICE_URL")
	if cartServiceUrl == "" {
		cartServiceUrl = ":8082"
	}
	lomsServiceUrl := os.Getenv("LOMS_SERVICE_URL")
	if lomsServiceUrl == "" {
		lomsServiceUrl = ":50777"
	}
	productServiceToken := os.Getenv("PRODUCT_SERVICE_TOKEN")
	if productServiceToken == "" {
		productServiceToken = "testtoken"
	}
	return Config{
		CartServiceUrl:      cartServiceUrl,
		LomsServiceUrl:      lomsServiceUrl,
		ProductServiceToken: productServiceToken,
	}
}
