package consumergroup

import (
	"context"
	"fmt"
	"route256/notifier/pkg/logger"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (cg *ConsumerGroup) Run(ctx context.Context) {
	logger.Infow(ctx, "[consumer-group] run")
	for {
		if err := cg.ConsumerGroup.Consume(ctx, cg.topics, cg.handler); err != nil {
			logger.Errorw(ctx, "[consumer-group] error from consume", "err", err)
		}
		if ctx.Err() != nil {
			logger.Infow(ctx, "[consumer-group] ctx closed", "err", ctx.Err().Error())
			return
		}
	}
}

func (cg *ConsumerGroup) RunErrorHandler(ctx context.Context) {
	for {
		select {
		case chErr, ok := <-cg.Errors():
			if !ok {
				logger.Infow(ctx, "[cg-error] error: chan closed")
				return
			}
			logger.Infow(ctx, "[cg-error] error", "err", chErr)
		case <-ctx.Done():
			logger.Infow(ctx, "[cg-error] ctx closed", "err", ctx.Err().Error())
			return
		}
	}
}

func NewConsumerGroup(brokers []string, groupID string, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...Option) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	
	for _, opt := range opts {
		if err := opt.Apply(config); err != nil {
			return nil, fmt.Errorf("opt.Apply: %w", err)
		}
	}
	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumerGroup: %w", err)
	}
	return &ConsumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}
