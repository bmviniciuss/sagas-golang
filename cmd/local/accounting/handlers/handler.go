package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/accounting/application"
	"github.com/bmviniciuss/sagas-golang/internal/streaming"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type AccountingKafkaHandler struct {
	logger      *zap.SugaredLogger
	publisher   streaming.Publisher
	handlersMap map[string]application.MessageHandler
}

var (
	_ streaming.Handler = (*AccountingKafkaHandler)(nil)
)

func NewAccountingKafkaHandler(logger *zap.SugaredLogger, publisher streaming.Publisher, handlersMap map[string]application.MessageHandler) *AccountingKafkaHandler {
	return &AccountingKafkaHandler{
		logger:      logger,
		publisher:   publisher,
		handlersMap: handlersMap,
	}
}

func (h *AccountingKafkaHandler) Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error {
	l := h.logger
	l.Infof("Accouting service received message [%s]", string(msg.Value))
	// TODO: add idempotence check

	var message events.Event
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		h.logger.With("error", err).Error("Got error unmarshalling message")
		return err
	}
	l.Infof("Successfully unmarshalled message")

	handler, ok := h.handlersMap[message.Type]
	if !ok {
		l.Infof("Ignoring message")
		err := commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}

	replyMessage, err := handler.Handle(ctx, &message)
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
	err = h.publisher.Publish(ctx, "service.accounting.events", data) // TODO: read from env
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
