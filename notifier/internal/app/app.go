package app

import (
	"context"
	"fmt"
	"route256/notifier/internal/pkg/config"
	consumergroup "route256/notifier/internal/pkg/infra/kafka/consumer_group"
	"sync"

	"github.com/IBM/sarama"
)

type App struct {
	consumergroup *consumergroup.ConsumerGroup
}

func NewApp(config config.Config) (*App, error) {
	cg, err := consumergroup.NewConsumerGroup(
		config.Brokers,
		config.ConsumerGroupName,
		[]string{config.OrderEventsTopic},
		consumergroup.NewConsumerGroupHandler(),
		consumergroup.WithOffsetsInitial(sarama.OffsetOldest),
		consumergroup.WithVersion(sarama.MaxVersion),
		consumergroup.WithConsumerReturnErrors(),
	)
	if err != nil {
		return nil, fmt.Errorf("consumergroup.NewConsumerGroup: %w", err)
	}
	return &App{
		consumergroup: cg,
	}, nil
}

func (app *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.consumergroup.Run(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.consumergroup.RunErrorHandler(ctx)
	}()
}

func (app *App) Stop() {
	app.consumergroup.Close()
}
