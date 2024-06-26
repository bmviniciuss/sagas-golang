package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/adapters/repositores/order" // TODO: fix typo
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/api"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/config/env"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/handlers"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
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

	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	lggr := logger.New(cfg.ServiceName)
	defer lggr.Sync()

	lggr.Info("Starting Order Service")
	dbpool, err := pgxpool.New(context.Background(), cfg.DBConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	lggr.Info("Connected to database")

	redisConn := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	err = redisConn.Ping(ctx).Err()
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	}

	var (
		bootstrapServers = cfg.KafkaBootstrapServers
		topics           = strings.Split(cfg.KafkaTopics, ",")
		group            = cfg.KafkaGroupID
	)

	publisher := streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	var (
		ordersRepository = order.NewRepositoryAdapter(lggr, dbpool)
		usecasesMap      = map[string]application.MessageHandler{
			"create_order":  handlers.NewCreateOrderHandler(lggr, usecases.NewCreateOrder(lggr, ordersRepository)),
			"approve_order": handlers.NewApproveOrder(lggr, usecases.NewApproveOrder(lggr, ordersRepository)),
			"reject_order":  handlers.NewRejectOrder(lggr, usecases.NewRejectOrder(lggr, ordersRepository)),
		}
	)

	handler := NewOrderMessageHandler(lggr, *publisher, usecasesMap)
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

	var (
		listUseCase  = usecases.NewListOrders(lggr, ordersRepository)
		getOrderByID = usecases.NewGetOrderByID(lggr, ordersRepository)
		apiHandlers  = api.NewHandlers(lggr, listUseCase, getOrderByID)
		httpServer   = newApiServer(":3000", apiHandlers)
	)

	go func() {
		lggr.Info("Starting Orders Service API go routine")
		if err := httpServer.ListenAndServe(); err != nil {
			lggr.With(zap.Error(err)).Fatal("Got error in Orders Service API server")
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

func newApiServer(addr string, handlers api.OrderHandlers) *http.Server {
	mux := api.NewRouter(handlers).Build()
	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
}
