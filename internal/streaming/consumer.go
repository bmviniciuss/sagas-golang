package streaming

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error
}

type Consumer struct {
	logger   *zap.SugaredLogger
	topics   []string
	consumer *kafka.Consumer
	running  bool
	handler  Handler
}

func NewConsumer(logger *zap.SugaredLogger, topics []string, kfkCfg *kafka.ConfigMap, handler Handler) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(kfkCfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		logger:   logger,
		topics:   topics,
		consumer: consumer,
		handler:  handler,
		running:  false,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) (err error) {
	l := c.logger
	l.Info("Starting consumer")
	err = c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error subscribing to topics")
		return err
	}
	defer func() {
		l.Info("Closing consumer")
		if err != nil {
			l.With(zap.Error(err)).Error("Closing with error")
		}
		cErr := c.consumer.Close()
		if cErr != nil {
			l.With(zap.Error(cErr)).Error("Got error closing consumer")
		}
	}()

	c.running = true
	for c.running {
		select {
		case <-ctx.Done():
			l.Info("Context done, stopping consumer")
			c.running = false
			return nil
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				l.Infof("Message received: %s", string(e.Value))
				err = c.handler.Handle(ctx, e, func() error {
					_, err := c.consumer.CommitMessage(e)
					return err
				})
				if err != nil {
					l.With(zap.Error(err)).Error("Got error handling message")
					return err
				}
			case kafka.Error:
				l.With(zap.Error(e)).Error("Got error")
			default:
				l.Infof("Ignored event: %v", e)
			}
		}
	}

	return nil
}
