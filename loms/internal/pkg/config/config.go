package config

import "os"

type Config struct {
	GrpcUrl string
	HttpUrl string
}

func NewConfig() Config {
	grpcUrl := os.Getenv("GRPC_URL")
	if grpcUrl == "" {
		grpcUrl = "localhost:50777"
	}
	httpUrl := os.Getenv("HTTP_URL")
	if httpUrl == "" {
		httpUrl = "localhost:8097"
	}
	return Config{
		GrpcUrl: grpcUrl,
		HttpUrl: httpUrl,
	}
}
