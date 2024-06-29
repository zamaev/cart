package config

import "os"

type Config struct {
	GrpcUrl      string
	HttpUrl      string
	DbMasterUrl  string
	DbReplicaUrl string
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
	dbMasterUrl := os.Getenv("DATABASE_MASTER_URL")
	if dbMasterUrl == "" {
		dbMasterUrl = "postgres://user:password@localhost:5432/postgres"
	}
	dbReplicaUrl := os.Getenv("DATABASE_REPLICA_URL")
	if dbReplicaUrl == "" {
		dbReplicaUrl = "postgres://user:password@localhost:5433/postgres"
	}
	return Config{
		GrpcUrl:      grpcUrl,
		HttpUrl:      httpUrl,
		DbMasterUrl:  dbMasterUrl,
		DbReplicaUrl: dbReplicaUrl,
	}
}
