package consumergroup

import "github.com/IBM/sarama"

type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

func WithVersion(v sarama.KafkaVersion) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Version = v
		return nil
	})
}

func WithOffsetsInitial(v int64) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Consumer.Offsets.Initial = v
		return nil
	})
}

func WithConsumerReturnErrors() Option {
	return optionFn(func(c *sarama.Config) error {
		c.Consumer.Return.Errors = true
		return nil
	})
}
