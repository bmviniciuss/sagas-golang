package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bmviniciuss/sagas-golang/cmd/local/customer/config/env"
	"github.com/bmviniciuss/sagas-golang/cmd/local/customer/handlers"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	lggr := logger.New(cfg.ServiceName)
	defer lggr.Sync()

	lggr.Info("Starting Customer Service")

	var (
		bootstrapServers = cfg.KafkaBootstrapServers
		topics           = strings.Split(cfg.KafkaTopics, ",")
		group            = cfg.KafkaGroupID
	)

	publisher := streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	verifyCustomerHandler := handlers.NewVerifyCustomer(lggr)
	handlersMap := map[string]application.MessageHandler{
		"verify_customer": verifyCustomerHandler,
	}
	handler := handlers.NewCustomerMessageHandler(lggr, *publisher, handlersMap)
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

	lggr.Info("Running Order Service. Waiting for signal to stop...")
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
