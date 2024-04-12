package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bmviniciuss/sagas-golang/internal/adapters/infra/kv"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CreateOrderPayloadBuilder struct {
}

func (b *CreateOrderPayloadBuilder) Build(ctx context.Context, data any, action saga.ActionType) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		sigCh       = make(chan os.Signal, 1)
		errCh       = make(chan error, 1)
	)
	defer cancel()
	defer close(sigCh)
	defer close(errCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	workflow := saga.Workflow{
		ID:           uuid.MustParse("2ef23373-9c01-4603-be2f-8e80552eb9a4"),
		Name:         "create_order",
		ReplyChannel: "saga.create_order.v1.response",
		Steps: saga.NewStepList(
			&saga.StepData{
				ID:             uuid.MustParse("4a4578ff-3602-4ad0-b262-6827c6ebc985"),
				Name:           "create_order",
				ServiceName:    "order",
				Compensable:    true,
				PayloadBuilder: &CreateOrderPayloadBuilder{},
			},
			&saga.StepData{
				ID:             uuid.MustParse("22d7c4bb-e751-4b47-a7a0-903ee5d3996e"),
				Name:           "verify_customer",
				ServiceName:    "customer",
				Compensable:    false,
				PayloadBuilder: &CreateOrderPayloadBuilder{},
			},
		),
	}
	lggr := logger.New()
	var (
		bootstrapServers = "localhost:9092"
		topics           = strings.Split("saga.create_order.v1.response", ",")
		group            = "sagas-golang"
	)

	redisConn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := redisConn.Ping(ctx).Err()
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	}

	idempotenceService := kv.NewAdapter(lggr, redisConn)
	publisher := streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	workflowService := service.NewWorkflow(lggr, publisher)
	messageHandler := streaming.NewMessageHandler(lggr, []saga.Workflow{workflow}, workflowService, idempotenceService)

	consumer, err := streaming.NewConsumer(lggr, topics, &kafka.ConfigMap{
		"bootstrap.servers":        bootstrapServers,
		"broker.address.family":    "v4",
		"group.id":                 group,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	}, messageHandler)
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error creating consumer")
	}

	go func() {
		if err := consumer.Start(ctx); err != nil {
			lggr.With(zap.Error(err)).Fatal("Got in consumer")
			errCh <- err
		}
	}()

	select {
	case <-sigCh:
		lggr.Info("Got signal, stopping consumer")
		cancel()
	case err := <-errCh:
		lggr.With(zap.Error(err)).Fatal("Got error in consumer")
		cancel()
		os.Exit(1)
	}

	lggr.Info("Exiting")
}
