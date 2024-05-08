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

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/adapters/repositores/ticket"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/api"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/usecases"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/handlers"
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
	lggr := logger.New("kitchen-service")
	defer lggr.Sync()

	lggr.Info("Starting Kitchen Service")
	dbpool, err := pgxpool.New(context.Background(), "postgres://sagas:sagas@localhost:5432/sagas") // TODO: add from env
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	lggr.Info("Connected to database")

	redisConn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // TODO: add from env
	})
	err = redisConn.Ping(ctx).Err()
	if err != nil {
		lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	}

	var (
		bootstrapServers = "localhost:9092"                              // TODO: add from env
		topics           = strings.Split("service.kitchen.request", ",") // TODO: add from env
		group            = "kitchen-service-group"                       // TODO: add from env
		publisher        = streaming.NewPublisher(lggr, &kafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		})
	)

	var (
		ticketRepository    = ticket.NewRepositoryAdapter(lggr, dbpool)
		createTicketUseCase = usecases.NewCreateTicket(lggr, ticketRepository)
		createTicketHandler = handlers.NewCreateTicketHandler(lggr, createTicketUseCase)
		handlersMap         = map[string]application.MessageHandler{
			"create_order_v1.create_ticket.request": createTicketHandler,
		}
		handler = handlers.NewTicketKafkaHandler(lggr, *publisher, handlersMap)
	)

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
		apiHandlers = api.NewHandlers(lggr)
		httpServer  = newApiServer(":3002", apiHandlers)
	)

	go func() {
		lggr.Info("Starting Kitchen Service API go routine")
		if err := httpServer.ListenAndServe(); err != nil {
			lggr.With(zap.Error(err)).Fatal("Got error in Kitchen Service API server")
			errCh <- err
		}
	}()

	lggr.Info("Running Kitchen Service. Waiting for signal to stop...")
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

func newApiServer(addr string, handlers api.KitchenApiHandler) *http.Server {
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
