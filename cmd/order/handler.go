package main

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type OrderMessageHandler struct {
	logger *zap.SugaredLogger
}

var (
	_ streaming.Handler = (*OrderMessageHandler)(nil)
)

func NewOrderMessageHandler(logger *zap.SugaredLogger) *OrderMessageHandler {
	return &OrderMessageHandler{
		logger: logger,
	}
}

func (h *OrderMessageHandler) Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error {
	h.logger.Infof("Handling message: %s", string(msg.Value))
	return nil
}
