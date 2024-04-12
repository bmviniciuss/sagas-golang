package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bmviniciuss/sagas-golang/internal/adapters/infra/kv"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

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

	lggr := logger.New()
	var (
		bootstrapServers = "localhost:9092"
		topics           = strings.Split("service.customer.request", ",")
		group            = "sagas-golang"
	)

	redisConn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := redisConn.Ping(ctx).Err()
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	}

	_ = kv.NewAdapter(lggr, redisConn)
	_ = streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	publisher := streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	handler := NewOrderMessageHandler(lggr, *publisher)
	consumer, err := streaming.NewConsumer(lggr, topics, &kafka.ConfigMap{
		"bootstrap.servers":        bootstrapServers,
		"broker.address.family":    "v4",
		"group.id":                 group,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	}, handler)
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
