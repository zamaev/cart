package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServiceName         string
	CartServiceUrl      string
	LomsServiceUrl      string
	ProductServiceUrl   string
	ProductServiceToken string
	TracerUrl           string
	RedisUrl            string
	RedisPassword       string
	RedisDB             int
	CacheSize           int
	CacheDefaultTTL     time.Duration
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
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "passwd"
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		cacheSize = 1000
	}
	cacheDefaultTTL, err := strconv.Atoi(os.Getenv("CACHE_DEFAULT_TTL"))
	if err != nil {
		cacheDefaultTTL = 60
	}
	return Config{
		ServiceName:         serviceName,
		CartServiceUrl:      cartServiceUrl,
		LomsServiceUrl:      lomsServiceUrl,
		ProductServiceUrl:   productServiceUrl,
		ProductServiceToken: productServiceToken,
		TracerUrl:           tracerUrl,
		RedisUrl:            redisUrl,
		RedisPassword:       redisPassword,
		RedisDB:             redisDB,
		CacheSize:           cacheSize,
		CacheDefaultTTL:     time.Duration(cacheDefaultTTL) * time.Second,
	}
}
