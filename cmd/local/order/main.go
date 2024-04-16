package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/adapters/repositores/order"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/api"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/internal/config/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	var (
		_, cancel = context.WithCancel(context.Background())
		sigCh     = make(chan os.Signal, 1)
		errCh     = make(chan error, 1)
	)
	defer cancel()
	defer close(sigCh)
	defer close(errCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	lggr := logger.New("orders-service")
	defer lggr.Sync()

	lggr.Info("Starting Order Service")
	dbpool, err := pgxpool.New(context.Background(), "postgres://sagas:sagas@localhost:5432/sagas")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	lggr.Info("Connected to database")

	// var (
	// 	bootstrapServers = "localhost:9092"
	// 	topics           = strings.Split("service.order.request", ",")
	// 	group            = "sagas-golang"
	// )

	// redisConn := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// })
	// err := redisConn.Ping(ctx).Err()
	// if err != nil {
	// 	lggr.With(zap.Error(err)).Fatal("Got error connecting to Redis")
	// }

	// _ = kv.NewAdapter(lggr, redisConn)
	// _ = streaming.NewPublisher(lggr, &kafka.ConfigMap{
	// 	"bootstrap.servers": bootstrapServers,
	// })
	// publisher := streaming.NewPublisher(lggr, &kafka.ConfigMap{
	// 	"bootstrap.servers": bootstrapServers,
	// })
	// handler := NewOrderMessageHandler(lggr, *publisher)
	// consumer, err := streaming.NewConsumer(lggr, topics, &kafka.ConfigMap{
	// 	"bootstrap.servers":        bootstrapServers,
	// 	"broker.address.family":    "v4",
	// 	"group.id":                 group,
	// 	"session.timeout.ms":       6000,
	// 	"auto.offset.reset":        "earliest",
	// 	"enable.auto.offset.store": false,
	// 	"enable.auto.commit":       true,
	// }, handler)
	// if err != nil {
	// 	lggr.With(zap.Error(err)).Fatal("Got error creating consumer")
	// }

	// go func() {
	// 	if err := consumer.Start(ctx); err != nil {
	// 		lggr.With(zap.Error(err)).Fatal("Got in consumer")
	// 		errCh <- err
	// 	}
	// }()

	var (
		ordersRepository = order.NewRepositoryAdapter(lggr, dbpool)
		listUseCase      = usecases.NewListOrders(lggr, ordersRepository)
		apiHandlers      = api.NewHandlers(lggr, listUseCase)
		httpServer       = newApiServer(":3001", apiHandlers)
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
