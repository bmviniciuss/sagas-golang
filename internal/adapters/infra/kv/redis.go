package kv

import (
	"context"
	"time"

	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Adapter struct {
	logger *zap.SugaredLogger
	client *redis.Client
}

func NewAdapter(
	logger *zap.SugaredLogger,
	client *redis.Client,
) *Adapter {
	return &Adapter{
		logger: logger,
		client: client,
	}
}

var (
	_ streaming.IdempotenceService = (*Adapter)(nil)
)

func (a *Adapter) Has(ctx context.Context, key string) (bool, error) {
	l := a.logger
	l.Infof("Checking if key [%s] exists", key)
	_, err := a.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			l.Infof("Key [%s] does not exist", key)
			return false, nil
		}
		l.With(zap.Error(err)).Error("Got error checking key")
		return false, err
	}
	return true, nil
}

func (a *Adapter) Set(ctx context.Context, key string, ttl time.Duration) error {
	l := a.logger
	l.Infof("Setting key [%s] with TTL [%s]", key, ttl)
	err := a.client.Set(ctx, key, "", ttl).Err()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error setting key")
		return err
	}
	return nil
}
