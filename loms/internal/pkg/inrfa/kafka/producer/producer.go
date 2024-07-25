package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/loms/internal/pkg/inrfa/kafka"
	"route256/loms/internal/pkg/model"
	"route256/loms/pkg/logger"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

type outboxRepository interface {
	Create(ctx context.Context, topic string, event []byte, headers []byte) (int64, error)
	GetWaitList(ctx context.Context) ([]model.OutboxItem, error)
	SetComplete(ctx context.Context, id int64) (err error)
}

type Producer struct {
	ctx                  context.Context
	close                context.CancelFunc
	producer             sarama.AsyncProducer
	handleEventsInterval time.Duration
	outboxRepository     outboxRepository
}

func (p *Producer) RunEventsHandle() {
	go p.runEventsResultHandle()

	ticker := time.NewTicker(p.handleEventsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			logger.Infow(p.ctx, "[handle events] terminate")
			return
		case <-ticker.C:
			logger.Infow(p.ctx, "[handle events] start")
			outboxItems, err := p.outboxRepository.GetWaitList(p.ctx)
			if err != nil {
				logger.Errorw(p.ctx, "[handle events] get wait list failed", "err", err)
				continue
			}

			for _, outboxItem := range outboxItems {
				var event model.Event
				if err := json.Unmarshal(outboxItem.Event, &event); err != nil {
					logger.Errorw(p.ctx, "[handle events] json.Unmarshal", "err", err, "id", outboxItem.Id)
					continue
				}

				var headers model.Headers
				if err := json.Unmarshal(outboxItem.Headers, &headers); err != nil {
					logger.Errorw(p.ctx, "[handle events] json.Unmarshal", "err", err, "id", outboxItem.Id)
					continue
				}

				msg := &sarama.ProducerMessage{
					Topic: outboxItem.Topic,
					Key:   sarama.StringEncoder(strconv.FormatInt(int64(event.OrderID), 10)),
					Value: sarama.ByteEncoder(outboxItem.Event),
					Headers: []sarama.RecordHeader{
						{
							Key:   []byte("TraceID"),
							Value: []byte(headers.TraceID),
						},
					},
				}

				p.producer.Input() <- msg
				if err := p.outboxRepository.SetComplete(p.ctx, outboxItem.Id); err != nil {
					logger.Errorw(p.ctx, "[handle events] set complete failed", "err", err)
				}
				logger.Infow(p.ctx, "[handle events] sent", "key", msg.Key, "topic", msg.Topic, "partition", msg.Partition)
			}
		}

	}
}

func (p *Producer) runEventsResultHandle() {
	for {
		select {
		case <-p.ctx.Done():
			logger.Infow(p.ctx, "[handle events result] terminate")
			return
		case msg := <-p.producer.Successes():
			if msg == nil {
				logger.Infow(p.ctx, "[handle events result] success chan closed")
				return
			}
			logger.Infow(p.ctx, "[handle events result] success", "key", msg.Key, "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
		case msgErr := <-p.producer.Errors():
			if msgErr == nil {
				logger.Infow(p.ctx, "[handle events result] error chan closed")
				return
			}
			logger.Infow(p.ctx, "[handle events result] error", "err", msgErr.Err, "topic", msgErr.Msg.Topic, "offset", msgErr.Msg.Offset)
		}
	}
}

func (p *Producer) Close() {
	p.close()
	p.producer.AsyncClose()
	<-p.producer.Successes()
	<-p.producer.Errors()
}

func NewProducer(ctx context.Context, kafkaConfig kafka.Config, outboxRepository outboxRepository, opts ...Option) (*Producer, error) {
	config := sarama.NewConfig()
	for _, opt := range opts {
		if err := opt.Apply(config); err != nil {
			return nil, fmt.Errorf("opt.Apply: %w", err)
		}
	}

	asyncProducer, err := sarama.NewAsyncProducer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	ctx, close := context.WithCancel(ctx)
	producer := &Producer{
		ctx:                  ctx,
		close:                close,
		producer:             asyncProducer,
		handleEventsInterval: time.Duration(kafkaConfig.HandleEventsInterval) * time.Second,
		outboxRepository:     outboxRepository,
	}

	return producer, nil
}
