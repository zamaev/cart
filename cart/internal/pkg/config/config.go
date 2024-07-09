package config

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName         string
	CartServiceUrl      string
	LomsServiceUrl      string
	ProductServiceUrl   string
	ProductServiceToken string
	TracerUrl           string
}

func NewConfig() Config {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "cart"
	}
	cartServiceUrl := os.Getenv("CART_SERVICE_URL")
	if cartServiceUrl == "" {
		cartServiceUrl = "localhost:8082"
	}
	lomsServiceUrl := os.Getenv("LOMS_SERVICE_URL")
	if lomsServiceUrl == "" {
		lomsServiceUrl = "localhost:50777"
	}
	productServiceUrl := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceUrl == "" {
		productServiceUrl = "http://route256.pavl.uk:8080"
	}
	productServiceToken := os.Getenv("PRODUCT_SERVICE_TOKEN")
	if productServiceToken == "" {
		productServiceToken = "testtoken"
	}
	tracerUrl := os.Getenv("TRACER_URL")
	if tracerUrl == "" {
		tracerUrl = "http://localhost:4318"
	}
	if !strings.HasPrefix(tracerUrl, "http") {
		tracerUrl = "http://" + tracerUrl
	}
	return Config{
		ServiceName:         serviceName,
		CartServiceUrl:      cartServiceUrl,
		LomsServiceUrl:      lomsServiceUrl,
		ProductServiceUrl:   productServiceUrl,
		ProductServiceToken: productServiceToken,
		TracerUrl:           tracerUrl,
	}
}
