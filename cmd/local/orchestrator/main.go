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

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/adapters/repositories/executions"
	workflowrepo "github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/adapters/repositories/workflows"
	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/api"
	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/config/env"
	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/workflows"
	"github.com/bmviniciuss/sagas-golang/internal/adapters/infra/kv"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/bmviniciuss/sagas-golang/pkg/validator"
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

	redisConn := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	err = redisConn.Ping(ctx).Err()
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DBConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	lggr.Info("Connected to database")

	var (
		workflows = []saga.Workflow{
			*workflows.NewCreateOrderV1(lggr),
		}
		workflowRepository = workflowrepo.NewInmemRepository(workflows)
	)

	var (
		executionsRepository = executions.NewRepositoryAdapter(lggr, dbpool, workflowRepository)
		bootstrapServers     = cfg.KafkaBootstrapServers
		topics               = strings.Split(cfg.KafkaTopics, ",")
		consumerGroupID      = cfg.KafkaGroupID
		publisher            = newPublisher(lggr, bootstrapServers)
		workflowService      = saga.NewService(lggr, executionsRepository, publisher)
		idempotenceService   = kv.NewAdapter(lggr, redisConn)
		messageHandler       = streaming.NewMessageHandler(lggr, executionsRepository, workflowService, idempotenceService)
	)

	var (
		val         = validator.New()
		apiHandlers = api.NewHandlers(lggr, workflowRepository, workflowService, val)
		httpServer  = newApiServer(":3000", apiHandlers)
	)

	consumer, err := newConsumer(lggr, topics, bootstrapServers, consumerGroupID, messageHandler)
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error creating consumer")
	}

	go func() {
		lggr.Infof("Starting orchestrator consumer go routine")
		if err := consumer.Start(ctx); err != nil {
			lggr.With(zap.Error(err)).Fatal("Got in consumer")
			errCh <- err
		}
	}()

	go func() {
		lggr.Info("Starting API server go routine")
		if err := httpServer.ListenAndServe(); err != nil {
			lggr.With(zap.Error(err)).Fatal("Got error in API server")
			errCh <- err
		}
	}()

	lggr.Info("Running. Waiting for signal to stop...")
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

func newApiServer(addr string, handlers api.HandlersPort) *http.Server {
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

func newPublisher(lggr *zap.SugaredLogger, servers string) *streaming.Publisher {
	return streaming.NewPublisher(lggr, &kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
}

func newConsumer(lggr *zap.SugaredLogger, topics []string, servers string, groupID string, handler streaming.Handler) (*streaming.Consumer, error) {
	return streaming.NewConsumer(lggr, topics, &kafka.ConfigMap{
		"bootstrap.servers":        servers,
		"broker.address.family":    "v4",
		"group.id":                 groupID,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	}, handler)
}
