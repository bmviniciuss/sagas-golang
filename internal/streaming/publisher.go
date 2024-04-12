package streaming

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Publisher struct { // TODO: add interface
	logger *zap.SugaredLogger
	kfkCfg *kafka.ConfigMap
}

func NewPublisher(logger *zap.SugaredLogger, kfkCfg *kafka.ConfigMap) *Publisher {
	return &Publisher{
		logger: logger,
		kfkCfg: kfkCfg,
	}
}

func (p *Publisher) Publish(ctx context.Context, destination string, data []byte) error {
	l := p.logger
	l.Infof("Publishing message to destination %s", destination)

	prod, err := kafka.NewProducer(p.kfkCfg)
	if err != nil {
		p.logger.With(zap.Error(err)).Error("Failed to create producer")
		return err
	}

	deliveryChan := make(chan kafka.Event, 100)
	defer close(deliveryChan)
	err = prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &destination, Partition: kafka.PartitionAny},
		Value:          data},
		deliveryChan,
	)
	if err != nil {
		p.logger.With(zap.Error(err)).Error("Failed to produce message")
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		p.logger.Errorf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		p.logger.Infof("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	return nil
}
