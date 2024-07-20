package config

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName  string
	GrpcUrl      string
	HttpUrl      string
	DbMasterUrl  string
	DbReplicaUrl string
	TracerUrl    string
}

func NewConfig() Config {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "loms"
	}
	grpcUrl := os.Getenv("GRPC_URL")
	if grpcUrl == "" {
		grpcUrl = "localhost:50777"
	}
	httpUrl := os.Getenv("HTTP_URL")
	if httpUrl == "" {
		httpUrl = "localhost:8097"
	}
	dbMasterUrl := os.Getenv("DATABASE_MASTER_URL")
	if dbMasterUrl == "" {
		dbMasterUrl = "postgres://user:password@localhost:5432/postgres"
	}
	dbReplicaUrl := os.Getenv("DATABASE_REPLICA_URL")
	if dbReplicaUrl == "" {
		dbReplicaUrl = "postgres://user:password@localhost:5433/postgres"
	}
	tracerUrl := os.Getenv("TRACER_URL")
	if tracerUrl == "" {
		tracerUrl = "http://localhost:4318"
	}
	if !strings.HasPrefix(tracerUrl, "http") {
		tracerUrl = "http://" + tracerUrl
	}
	return Config{
		ServiceName:  serviceName,
		GrpcUrl:      grpcUrl,
		HttpUrl:      httpUrl,
		DbMasterUrl:  dbMasterUrl,
		DbReplicaUrl: dbReplicaUrl,
		TracerUrl:    tracerUrl,
	}
}
