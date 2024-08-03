package config

import (
	"fmt"
	"os"
	"route256/loms/internal/pkg/inrfa/kafka"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName       string
	GrpcUrl           string
	HttpUrl           string
	DbMasterUrl       string
	DbReplicaUrl      string
	DbMasterShard2Url string
	TracerUrl         string
	Kafka             kafka.Config
}

func NewConfig() (Config, error) {
	var err error

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
	dbMasterShard2Url := os.Getenv("DATABASE_MASTER-SHARD-2_URL")
	if dbMasterShard2Url == "" {
		dbMasterShard2Url = "postgres://user:password@localhost:5434/postgres"
	}

	tracerUrl := os.Getenv("TRACER_URL")
	if tracerUrl == "" {
		tracerUrl = "http://localhost:4318"
	}
	if !strings.HasPrefix(tracerUrl, "http") {
		tracerUrl = "http://" + tracerUrl
	}

	kafkaBrokers := []string{"localhost:9092"}
	kafkaBrokersRaw := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersRaw != "" {
		kafkaBrokers = strings.Split(kafkaBrokersRaw, ",")
	}
	kafkaOrderEventsTopic := os.Getenv("KAFKA_ORDER_EVENTS_TOPIC")
	if kafkaOrderEventsTopic == "" {
		kafkaOrderEventsTopic = "loms.order-events"
	}
	handleEventsInterval := int64(5)
	handleEventsIntervalRaw := os.Getenv("HANDLE_EVENTS_INTERVAL")
	if handleEventsIntervalRaw != "" {
		handleEventsInterval, err = strconv.ParseInt(handleEventsIntervalRaw, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
		}
	}

	return Config{
		ServiceName:       serviceName,
		GrpcUrl:           grpcUrl,
		HttpUrl:           httpUrl,
		DbMasterUrl:       dbMasterUrl,
		DbReplicaUrl:      dbReplicaUrl,
		DbMasterShard2Url: dbMasterShard2Url,
		TracerUrl:         tracerUrl,
		Kafka: kafka.Config{
			Brokers:              kafkaBrokers,
			OrderEventsTopic:     kafkaOrderEventsTopic,
			HandleEventsInterval: handleEventsInterval,
		},
	}, nil
}
