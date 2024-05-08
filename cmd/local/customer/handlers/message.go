package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type CustomerMessageHandler struct {
	logger      *zap.SugaredLogger
	publisher   streaming.Publisher
	handlersMap map[string]application.MessageHandler
}

var (
	_ streaming.Handler = (*CustomerMessageHandler)(nil)
)

func NewCustomerMessageHandler(logger *zap.SugaredLogger, publisher streaming.Publisher, handlersMap map[string]application.MessageHandler) *CustomerMessageHandler {
	return &CustomerMessageHandler{
		logger:      logger,
		publisher:   publisher,
		handlersMap: handlersMap,
	}
}

func (h *CustomerMessageHandler) Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error {
	l := h.logger
	l.Infof("Customer service received message [%s]", string(msg.Value))
	// TODO: add idempotence check

	var message saga.Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		h.logger.With("error", err).Error("Got error unmarshalling message")
		return err
	}
	l.Infof("Successfully unmarshalled message")

	useCase, ok := h.handlersMap[message.EventType.String()]
	if !ok {
		l.Infof("Ignoring message")
		err := commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}

	replyMessage, err := useCase.Handle(ctx, &message)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error handling message")
		return err
	}

	data, err := json.Marshal(replyMessage)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error marshalling reply message")
		return err
	}
	l.Infof("Successfully marshalled reply message")
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
