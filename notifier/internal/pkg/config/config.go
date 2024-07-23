package config

import (
	"os"
	"strings"
)

type Config struct {
	Brokers           []string
	ConsumerGroupName string
	OrderEventsTopic  string
}

func NewConfig() Config {
	brokers := []string{"localhost:9092"}
	brokersRaw := os.Getenv("KAFKA_BROKERS")
	if brokersRaw != "" {
		brokers = strings.Split(brokersRaw, ",")
	}
	consumerGroupName := os.Getenv("KAFKA_CONSUMER_GROUP_NAME")
	if consumerGroupName == "" {
		consumerGroupName = "notifier"
	}
	orderEventsTopic := os.Getenv("KAFKA_ORDER_EVENTS_TOPIC")
	if orderEventsTopic == "" {
		orderEventsTopic = "loms.order-events"
	}
	return Config{
		Brokers:           brokers,
		ConsumerGroupName: consumerGroupName,
		OrderEventsTopic:  orderEventsTopic,
	}
}
