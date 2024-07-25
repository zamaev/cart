package consumergroup

import (
	"encoding/json"
	"route256/notifier/internal/pkg/model"
	"route256/notifier/pkg/logger"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct{}

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

func NewConsumerGroupHandler() *ConsumerGroupHandler {
	return &ConsumerGroupHandler{}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event model.Event
		json.Unmarshal(msg.Value, &event)

		logger.Infow(session.Context(), "[kafka-consumer] consume event", "event", event)

		session.MarkMessage(msg, "")
	}
	return nil
}
