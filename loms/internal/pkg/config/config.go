package config

import "os"

type Config struct {
	GrpcUrl string
	HttpUrl string
}

func NewConfig() Config {
	grpcUrl := os.Getenv("GRPC_URL")
	if grpcUrl == "" {
		grpcUrl = ":50777"
	}
	httpUrl := os.Getenv("HTTP_URL")
	if httpUrl == "" {
		httpUrl = ":8097"
	}
	return Config{
		GrpcUrl: grpcUrl,
		HttpUrl: httpUrl,
	}
}
