package config

import "os"

type Config struct {
	ProductServiceToken string
	ProductServiceUrl   string
}

func NewConfig() Config {
	productServiceToken := os.Getenv("PRODUCT_SERVICE_TOKEN")
	if productServiceToken == "" {
		productServiceToken = "testtoken"
	}
	return Config{
		ProductServiceToken: productServiceToken,
	}
}
