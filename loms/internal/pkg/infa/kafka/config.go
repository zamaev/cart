package kafka

type Config struct {
	Brokers              []string
	OrderEventsTopic     string
	HandleEventsInterval int64
}
