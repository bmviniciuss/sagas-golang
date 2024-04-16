package main

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type OrderMessageHandler struct {
	logger    *zap.SugaredLogger
	publisher streaming.Publisher
}

var (
	_ streaming.Handler = (*OrderMessageHandler)(nil)
)

func NewOrderMessageHandler(logger *zap.SugaredLogger, publisher streaming.Publisher) *OrderMessageHandler {
	return &OrderMessageHandler{
		logger:    logger,
		publisher: publisher,
	}
}

func (h *OrderMessageHandler) Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error {
	l := h.logger
	l.Infof("Order service received message [%s]", string(msg.Value))

	var message saga.Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		h.logger.With("error", err).Error("Got error unmarshalling message")
		return err
	}
	l.Infof("Successfully unmarshalled message")

	if message.EventType.SagaName != "create_order" && message.EventType.StepName != "create_order" {
		l.Infof("Ignoring message")
		err := commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}
	l.Infof("Handling message create_order.create_order message")
	replyMessage := saga.NewParticipantMessage(message.GlobalID, nil, nil, saga.SuccessActionType, &message)
	data, err := json.Marshal(replyMessage)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error marshalling reply message")
		return err
	}

	err = h.publisher.Publish(ctx, message.Saga.ReplyChannel, data)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error publishing message")
		return err
	}

	err = commitFn()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error committing message")
		return err
	}

	return nil
}
